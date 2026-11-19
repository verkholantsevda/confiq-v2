package backup

import "errors"

var (
	ErrBackupNotFound       = errors.New("backup not found")
	ErrInvalidBackupName    = errors.New("invalid backup name")
	ErrInvalidBackup        = errors.New("invalid backup file")
	ErrBackupDirUnavailable = errors.New("backup directory is not available")
)
