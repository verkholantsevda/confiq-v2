package logger

import "log/slog"

func audit(action string, args ...any) {
	fields := append([]any{"action", action}, args...)
	slog.Info("audit", fields...)
}

func UserLogin(id uint, username string) {
	audit(
		"user.login",
		"user_id", id,
		"username", username,
	)
}

func UserCreated(id uint, username string) {
	audit(
		"user.created",
		"user_id", id,
		"username", username,
	)
}

func UserUpdated(id uint, username string) {
	audit(
		"user.updated",
		"user_id", id,
		"username", username,
	)
}

func UserDeleted(id uint, username string) {
	audit(
		"user.deleted",
		"user_id", id,
		"username", username,
	)
}

func PasswordChanged(id uint, username string) {
	audit(
		"user.password_changed",
		"user_id", id,
		"username", username,
	)
}

func ConfigCreated(userID, configID uint, name string) {
	audit(
		"config.created",
		"user_id", userID,
		"config_id", configID,
		"name", name,
	)
}

func ConfigUpdated(userID, configID uint, name string) {
	audit(
		"config.created",
		"user_id", userID,
		"config_id", configID,
		"name", name,
	)
}

func ConfigDeleted(userID, configID uint) {
	audit(
		"config.deleted",
		"user_id", userID,
		"config_id", configID,
	)
}
