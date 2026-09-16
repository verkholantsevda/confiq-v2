package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"confiq/internal/audit"
	"confiq/internal/auth"
	"confiq/internal/config"
	"confiq/internal/configs"
	"confiq/internal/configtypes"
	"confiq/internal/database"
	"confiq/internal/endpoints"
	"confiq/internal/groups"
	"confiq/internal/logger"
	"confiq/internal/middleware"
	"confiq/internal/password"
	"confiq/internal/router"
	"confiq/internal/users"
	"confiq/internal/warp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}
	logger.Init(cfg.LogLevel)
	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB from gorm", "error", err)
		return
	}

	if err := database.RunMigrations(sqlDB); err != nil {
		slog.Error("failed to apply migrations", "error", err)
		return
	}
	// Repositories
	userRepo := users.NewRepository(db)
	groupRepo := groups.NewRepository(db)
	endpointRepo := endpoints.NewRepository(db)
	configTypeRepo := configtypes.NewRepository(db)
	configRepo := configs.NewRepository(db)
	auditRepo := audit.NewRepository(db)
	auditService := audit.NewService(auditRepo)

	// Services
	userService := users.NewService(userRepo, cfg, auditService)
	groupService := groups.NewService(
		groupRepo,
		endpointRepo,
		auditService,
	)
	endpointService := endpoints.NewService(endpointRepo, auditService)
	configTypeService := configtypes.NewService(configTypeRepo, auditService)
	generator := warp.NewGenerator()
	configService := configs.NewService(
		configRepo,
		userRepo,
		endpointRepo,
		configTypeRepo,
		generator,
		auditService,
	)
	hash, err := password.HashPassword(cfg.AdminPassword)
	if err != nil {
		slog.Error("failed to hash admin password", "error", err)
		return
	}

	if err := userService.EnsureAdmin(cfg.AdminUsername, hash); err != nil {
		slog.Error("failed to ensure admin user", "error", err)
		return
	}
	jwtService := auth.NewJWT(cfg.JWTSecret, cfg.JWTExpire)
	authService := auth.NewService(userService, jwtService, auditService)

	// Handlers
	userHandler := users.NewHandler(userService)
	groupHandler := groups.NewHandler(groupService)
	endpointHandler := endpoints.NewHandler(endpointService)
	configTypeHandler := configtypes.NewHandler(configTypeService)
	configHandler := configs.NewHandler(configService)
	authHandler := auth.NewHandler(authService)
	auditHandler := audit.NewHandler(auditService)

	authMiddleware := middleware.NewAuth(jwtService)

	// Router
	r := router.New(userHandler, groupHandler, endpointHandler, configHandler, configTypeHandler, authHandler, auditHandler, authMiddleware)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: r,
	}
	slog.Info(
		"starting application",
		"name", cfg.AppName,
		"version", cfg.AppVersion,
		"env", cfg.AppEnv,
		"address", server.Addr,
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server stopped", "error", err)
	}
}
