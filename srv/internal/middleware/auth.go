package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"regres/srv/internal/auth"
	"regres/srv/internal/database/queries"

	"github.com/jackc/pgx/v5"
)

type contextKey string

const currentUserKey contextKey = "currentUser"

type CurrentUser struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type AuthMiddleware struct {
	Queries *queries.Queries
}

func NewAuthMiddleware(q *queries.Queries) *AuthMiddleware {
	return &AuthMiddleware{
		Queries: q,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := auth.ReadSessionCookie(r)
		if !ok {
			writeAuthError(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		tokenHash := auth.HashSessionToken(token)

		session, err := m.Queries.GetSessionByTokenHash(r.Context(), tokenHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writeAuthError(w, http.StatusUnauthorized, "not authenticated")
				return
			}

			writeAuthError(w, http.StatusInternalServerError, "could not load session")
			return
		}

		_ = m.Queries.UpdateSessionLastSeen(r.Context(), session.ID)

		user := CurrentUser{
			ID:       session.UserID,
			Email:    session.Email,
			Username: session.Username,
		}

		ctx := context.WithValue(r.Context(), currentUserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	user, ok := ctx.Value(currentUserKey).(CurrentUser)
	return user, ok
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
