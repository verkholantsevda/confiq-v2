package backup

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"confiq/internal/audit"
	"confiq/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var backupNamePattern = regexp.MustCompile(`^confiq-\d{8}-\d{6}(-\d+)?\.db$`)

type Service struct {
	cfg      *config.Config
	db       *gorm.DB
	audit    *audit.Service
	shutdown func()

	mu sync.Mutex
}

func NewService(cfg *config.Config, db *gorm.DB, auditService *audit.Service, shutdown func()) (*Service, error) {
	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBackupDirUnavailable, err)
	}

	return &Service{
		cfg:      cfg,
		db:       db,
		audit:    auditService,
		shutdown: shutdown,
	}, nil
}

// StartScheduler создаёт начальный бэкап, если его ещё нет, и дальше делает
// копии с заданным интервалом. Ручные бэкапы продолжают работать и при
// BACKUP_ENABLED=false.
func (s *Service) StartScheduler(ctx context.Context) {
	if !s.cfg.BackupEnabled {
		return
	}

	go func() {
		files, err := s.List()
		if err == nil && len(files) == 0 {
			if _, err := s.Create(); err != nil {
				slog.Error("initial database backup failed", "error", err)
			}
		}

		ticker := time.NewTicker(s.cfg.BackupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := s.Create(); err != nil {
					slog.Error("scheduled database backup failed", "error", err)
				}
			}
		}
	}()
}

func (s *Service) List() ([]FileInfo, error) {
	entries, err := os.ReadDir(s.cfg.BackupDir)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || !backupNamePattern.MatchString(entry.Name()) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:      entry.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].CreatedAt.After(files[j].CreatedAt)
	})

	return files, nil
}

func (s *Service) Create() (FileInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.createLocked()
}

func (s *Service) createLocked() (FileInfo, error) {
	name, path, err := s.nextPath()
	if err != nil {
		return FileInfo{}, err
	}

	// VACUUM INTO даёт консистентный снимок, не мешая активным подключениям.
	statement := fmt.Sprintf("VACUUM INTO '%s'", escapeSQLString(path))
	if err := s.db.Exec(statement).Error; err != nil {
		return FileInfo{}, err
	}

	if err := s.applyRetentionLocked(); err != nil {
		slog.Error("failed to prune old database backups", "error", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Name:      name,
		Size:      info.Size(),
		CreatedAt: info.ModTime(),
	}, nil
}

func (s *Service) Path(name string) (string, error) {
	if !backupNamePattern.MatchString(name) {
		return "", ErrInvalidBackupName
	}

	return filepath.Join(s.cfg.BackupDir, name), nil
}

func (s *Service) Delete(name string) error {
	path, err := s.Path(name)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return ErrBackupNotFound
		}

		return err
	}

	return nil
}

// CreateTempFile создаёт временный файл рядом с БД: в таком случае restore
// делается обычным rename в пределах одной файловой системы, даже когда
// каталог бэкапов смонтирован отдельным volume.
func (s *Service) CreateTempFile() (*os.File, error) {
	dir := filepath.Dir(s.cfg.DatabasePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return os.CreateTemp(dir, ".restore-*.db")
}

// StageBackup копирует уже существующий бэкап рядом с БД. BACKUP_DIR может
// быть отдельным монтом, поэтому напрямую rename оттуда в DATABASE_PATH
// нельзя — сначала файл попадает во временный каталог рядом с БД.
func (s *Service) StageBackup(name string) (string, error) {
	path, err := s.Path(name)
	if err != nil {
		return "", err
	}

	source, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrBackupNotFound
		}

		return "", err
	}
	defer source.Close()

	temp, err := s.CreateTempFile()
	if err != nil {
		return "", err
	}

	tempPath := temp.Name()

	if _, err := io.Copy(temp, source); err != nil {
		temp.Close()
		os.Remove(tempPath)

		return "", err
	}

	if err := temp.Close(); err != nil {
		os.Remove(tempPath)

		return "", err
	}

	return tempPath, nil
}

func (s *Service) Validate(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make([]byte, 16)
	if _, err := file.Read(header); err != nil {
		return fmt.Errorf("%w: file is too short", ErrInvalidBackup)
	}

	if string(header) != "SQLite format 3\x00" {
		return fmt.Errorf("%w: not a SQLite database", ErrInvalidBackup)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBackup, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var integrity string
	if err := sqlDB.QueryRow("PRAGMA integrity_check").Scan(&integrity); err != nil {
		return fmt.Errorf("%w: integrity check failed: %v", ErrInvalidBackup, err)
	}

	if !strings.EqualFold(integrity, "ok") {
		return fmt.Errorf("%w: %s", ErrInvalidBackup, integrity)
	}

	var tables int
	if err := sqlDB.QueryRow(
		"SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name IN ('users', 'configs')",
	).Scan(&tables); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBackup, err)
	}

	if tables < 2 {
		return fmt.Errorf("%w: database schema is not recognized", ErrInvalidBackup)
	}

	return nil
}

// Restore при BACKUP_SAFETY_ON_RESTORE=true сохраняет текущую БД в бэкап,
// затем подменяет файл. Вызывающий должен завершить процесс, чтобы GORM
// переоткрыл новую БД.
func (s *Service) Restore(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cfg.BackupSafetyOnRestore {
		if _, err := s.createLocked(); err != nil {
			return fmt.Errorf("safety backup failed: %w", err)
		}
	}

	if err := os.Rename(path, s.cfg.DatabasePath); err != nil {
		return err
	}

	return nil
}

func (s *Service) Shutdown() {
	if s.shutdown != nil {
		s.shutdown()
	}
}

func (s *Service) Log(userID *uint, action string, message string) {
	if s.audit == nil {
		return
	}

	if err := s.audit.Log(userID, action, "backup", nil, message); err != nil {
		slog.Error("failed to write audit log", "error", err)
	}
}

// LogTo записывает событие прямо в восстанавливаемый файл: после подмены
// основная БД уже недоступна для записи старым подключением (SQLite видит
// смену inode и возвращает SQLITE_READONLY_DBMOVED).
func (s *Service) LogTo(path string, userID *uint, action string, message string) error {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	return db.Create(&audit.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: "backup",
		Message:    message,
		CreatedAt:  time.Now(),
	}).Error
}

func (s *Service) nextPath() (string, string, error) {
	timestamp := time.Now().UTC().Format("20060102-150405")
	name := fmt.Sprintf("confiq-%s.db", timestamp)
	path := filepath.Join(s.cfg.BackupDir, name)

	for suffix := 1; ; suffix++ {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return name, path, nil
		}

		name = fmt.Sprintf("confiq-%s-%d.db", timestamp, suffix)
		path = filepath.Join(s.cfg.BackupDir, name)
	}
}

func (s *Service) applyRetentionLocked() error {
	if s.cfg.BackupKeep <= 0 {
		return nil
	}

	files, err := s.List()
	if err != nil {
		return err
	}

	for i := s.cfg.BackupKeep; i < len(files); i++ {
		if err := os.Remove(filepath.Join(s.cfg.BackupDir, files[i].Name)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func escapeSQLString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
