package backup

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"confiq/internal/httpx"
	"confiq/internal/identity"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	files, err := h.service.List()
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	httpx.OK(w, files)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := identity.GetClaims(r)

	info, err := h.service.Create()
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	if claims != nil {
		h.service.Log(&claims.UserID, "backup.created", "Создан бэкап базы данных "+info.Name)
	}

	httpx.Created(w, info)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	path, err := h.service.Path(name)
	if err != nil {
		httpx.BadRequest(w, err.Error())
		return
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			httpx.NotFound(w, ErrBackupNotFound.Error())
			return
		}

		httpx.InternalServerError(w, err)
		return
	}

	claims := identity.GetClaims(r)
	if claims != nil {
		h.service.Log(&claims.UserID, "backup.downloaded", "Скачан бэкап базы данных "+name)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	http.ServeFile(w, r, path)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := h.service.Delete(name); err != nil {
		switch err {
		case ErrBackupNotFound:
			httpx.NotFound(w, err.Error())
		case ErrInvalidBackupName:
			httpx.BadRequest(w, err.Error())
		default:
			httpx.InternalServerError(w, err)
		}

		return
	}

	claims := identity.GetClaims(r)
	if claims != nil {
		h.service.Log(&claims.UserID, "backup.deleted", "Удалён бэкап базы данных "+name)
	}

	httpx.NoContent(w)
}

func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		httpx.BadRequest(w, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.BadRequest(w, "file is required")
		return
	}
	defer file.Close()

	temp, err := h.service.CreateTempFile()
	if err != nil {
		httpx.InternalServerError(w, err)
		return
	}

	tempPath := temp.Name()

	if _, err := io.Copy(temp, file); err != nil {
		temp.Close()
		os.Remove(tempPath)
		httpx.InternalServerError(w, err)
		return
	}

	if err := temp.Close(); err != nil {
		os.Remove(tempPath)
		httpx.InternalServerError(w, err)
		return
	}

	h.restoreStaged(w, r, tempPath, header.Filename)
}

// RestoreExisting восстанавливает БД из бэкапа, который уже лежит на сервере.
func (h *Handler) RestoreExisting(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	tempPath, err := h.service.StageBackup(name)
	if err != nil {
		switch err {
		case ErrBackupNotFound:
			httpx.NotFound(w, err.Error())
		case ErrInvalidBackupName:
			httpx.BadRequest(w, err.Error())
		default:
			httpx.InternalServerError(w, err)
		}

		return
	}

	h.restoreStaged(w, r, tempPath, name)
}

func (h *Handler) restoreStaged(w http.ResponseWriter, r *http.Request, tempPath string, sourceName string) {
	if err := h.service.Validate(tempPath); err != nil {
		os.Remove(tempPath)
		httpx.BadRequest(w, err.Error())
		return
	}

	claims := identity.GetClaims(r)
	if claims != nil {
		if err := h.service.LogTo(
			tempPath,
			&claims.UserID,
			"backup.restored",
			"Восстановлена база данных из бэкапа "+sourceName,
		); err != nil {
			slog.Error("failed to write audit log into restored database", "error", err)
		}
	}

	if err := h.service.Restore(tempPath); err != nil {
		os.Remove(tempPath)
		httpx.InternalServerError(w, err)
		return
	}

	httpx.OK(w, map[string]bool{"restarting": true})

	go func() {
		time.Sleep(300 * time.Millisecond)
		h.service.Shutdown()
	}()
}
