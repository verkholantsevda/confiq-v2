package router

import (
	"net/http"
	"os"

	"confiq/internal/audit"
	"confiq/internal/auth"
	"confiq/internal/configs"
	"confiq/internal/configtypes"
	"confiq/internal/endpoints"
	"confiq/internal/groups"
	authmw "confiq/internal/middleware"
	"confiq/internal/users"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(userHandler *users.Handler, groupHandler *groups.Handler, endpointHandler *endpoints.Handler, configHandler *configs.Handler, configTypeHandler *configtypes.Handler, authHandler *auth.Handler, auditHandler *audit.Handler, authMiddleware *authmw.Auth) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost"
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			allowedOrigin,
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},

		ExposedHeaders: []string{
			"Link",
		},

		AllowCredentials: true,

		MaxAge: 300,
	}))

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API
	r.Route("/api", func(r chi.Router) {

		r.Post("/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticated)

			r.Route("/me", func(r chi.Router) {
				r.Get("/", userHandler.Me)
				r.Put("/password", userHandler.ChangePassword)
			})
			r.Route("/totp", func(r chi.Router) {
				r.Get("/status", userHandler.TOTPStatus)
				r.Post("/setup", userHandler.SetupTOTP)
				r.Post("/enable", userHandler.EnableTOTP)
				r.Delete("/", userHandler.DisableTOTP)
			})
			r.Route("/configs", func(r chi.Router) {
				r.Get("/", configHandler.List)
				r.Get("/{id}", configHandler.Get)
				r.Post("/", configHandler.Create)
				r.Put("/{id}", configHandler.Update)
				r.Delete("/{id}", configHandler.Delete)
			})
			r.Route("/endpoints", func(r chi.Router) {
				r.Get("/", endpointHandler.ListForUser)
				r.Get("/{id}", endpointHandler.Get)
				r.Get("/{id}/config-types", endpointHandler.ListConfigTypes)
			})
			r.Route("/config-types-user", func(r chi.Router) {
				r.Get("/", configTypeHandler.List)
			})

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAdmin)

				r.Get("/activity", auditHandler.List)

				r.Route("/users", func(r chi.Router) {
					r.Get("/", userHandler.List)
					r.Post("/", userHandler.Create)
					r.Put("/{id}", userHandler.Update)
					r.Delete("/{id}", userHandler.Delete)
				})

				r.Route("/groups", func(r chi.Router) {
					r.Get("/", groupHandler.List)
					r.Get("/{id}", groupHandler.Get)
					r.Post("/", groupHandler.Create)
					r.Put("/{id}", groupHandler.Update)
					r.Delete("/{id}", groupHandler.Delete)
				})

				r.Route("/endpoints-all", func(r chi.Router) {
					r.Get("/", endpointHandler.List)
					r.Get("/{id}", endpointHandler.Get)
					r.Post("/", endpointHandler.Create)
					r.Put("/{id}", endpointHandler.Update)
					r.Delete("/{id}", endpointHandler.Delete)
				})

				r.Route("/groups_endpoints", func(r chi.Router) {
					r.Get("/", endpointHandler.ListGroupsEndpoints)
					r.Post("/", endpointHandler.AddGroupEndpoint)
					r.Delete("/{group_id}/{endpoint_id}", endpointHandler.RemoveGroupEndpoint)
				})

				r.Route("/config-types", func(r chi.Router) {
					r.Get("/", configTypeHandler.List)
					r.Get("/{id}", configTypeHandler.Get)
					r.Post("/", configTypeHandler.Create)
					r.Put("/{id}", configTypeHandler.Update)
					r.Delete("/{id}", configTypeHandler.Delete)
				})

				r.Get("/all-configs", configHandler.ListAll)
			})
		})
	})

	return r
}
