package routes

import (
	"net/http"

	"regres/srv/internal/database/queries"
	"regres/srv/internal/handlers"
	"regres/srv/internal/middleware"
)

func Register(mux *http.ServeMux, q *queries.Queries) {
	authHandler := handlers.NewAuthHandler(q)
	authMiddleware := middleware.NewAuthMiddleware(q)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Public routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Protected routes
	mux.Handle("POST /api/auth/logout", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))

	// Example protected route
	mux.Handle("GET /api/protected", authMiddleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.CurrentUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handlers.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "you are authenticated",
			"user":    user,
		})
	})))
}
