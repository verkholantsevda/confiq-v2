package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	AppVersion    string
	AppEnv        string
	AdminUsername string
	AdminPassword string
	Host          string
	Port          string

	DatabasePath string

	JWTSecret string
	JWTExpire time.Duration

	CFAPI     string
	CFTimeout time.Duration

	LogLevel string

	EnableSwagger bool

	OTPUserEnabled  bool
	OTPAdminEnabled bool
}

func Load() (*Config, error) {
	// Если .env отсутствует — ничего страшного.
	_ = godotenv.Load()

	timeout, err := strconv.Atoi(getEnv("CF_TIMEOUT", "30"))
	if err != nil {
		timeout = 30
	}

	jwtExpire, err := time.ParseDuration(getEnv("JWT_EXPIRE", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRE: %w", err)
	}

	cfg := &Config{
		AppName:       getEnv("APP_NAME", "Confiq"),
		AppVersion:    getEnv("APP_VERSION", "0.1.0"),
		AppEnv:        getEnv("APP_ENV", "development"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin"),
		Host:          getEnv("HOST", "0.0.0.0"),
		Port:          getEnv("PORT", "8080"),

		DatabasePath: getEnv("DATABASE_PATH", "data/confiq.db"),

		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpire: jwtExpire,

		CFAPI:     getEnv("CF_API", "https://api.cloudflareclient.com/v0i1909051800"),
		CFTimeout: time.Duration(timeout) * time.Second,

		LogLevel: getEnv("LOG_LEVEL", "info"),

		EnableSwagger:   getEnv("ENABLE_SWAGGER", "false") == "true",
		OTPUserEnabled:  getEnv("OTP_USER", "false") == "true",
		OTPAdminEnabled: getEnv("OTP_ADMIN", "true") == "true",
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
