package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"regres/srv/internal/auth"
	"regres/srv/internal/config"
	"regres/srv/internal/database/queries"
	"regres/srv/internal/helpers"
	"regres/srv/internal/middleware"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuthHandler struct {
	Queries *queries.Queries
}

func NewAuthHandler(q *queries.Queries) *AuthHandler {
	return &AuthHandler{
		Queries: q,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type AuthResponse struct {
	User UserResponse `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	body.Username = strings.TrimSpace(body.Username)

	if body.Email == "" || body.Username == "" || body.Password == "" {
		helpers.WriteError(w, http.StatusBadRequest, "email, username, and password are required")
		return
	}

	if len(body.Password) < 8 {
		helpers.WriteError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	passwordHash, err := auth.HashPassword(body.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := h.Queries.CreateUser(r.Context(), queries.CreateUserParams{
		Email:        body.Email,
		Username:     body.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "could not create user")
		return
	}

	sessionToken, expiresAt, err := h.createSession(r.Context(), r, user.ID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	auth.SetSessionCookie(w, sessionToken, expiresAt)

	helpers.WriteJSON(w, http.StatusCreated, AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))

	user, err := h.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if !auth.CheckPassword(body.Password, user.PasswordHash) {
		helpers.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	sessionToken, expiresAt, err := h.createSession(r.Context(), r, user.ID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	auth.SetSessionCookie(w, sessionToken, expiresAt)

	helpers.WriteJSON(w, http.StatusOK, AuthResponse{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := auth.ReadSessionCookie(r)
	if ok {
		tokenHash := auth.HashSessionToken(token)
		_ = h.Queries.DeleteSession(r.Context(), tokenHash)
	}

	auth.ClearSessionCookie(w)

	helpers.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.CurrentUserFromContext(r.Context())
	if !ok {
		helpers.WriteError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (h *AuthHandler) createSession(ctx context.Context, r *http.Request, userID int64) (string, time.Time, error) {
	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		return "", time.Time{}, err
	}

	tokenHash := auth.HashSessionToken(sessionToken)
	expiresAt := time.Now().Add(time.Duration(config.SessionDurationHours()) * time.Hour)

	_, err = h.Queries.CreateSession(ctx, queries.CreateSessionParams{
		UserID:           userID,
		SessionTokenHash: tokenHash,
		UserAgent:        pgText(r.UserAgent()),
		IpAddress:        pgText(getIP(r)),
		ExpiresAt:        pgTimestamptz(expiresAt),
	})
	if err != nil {
		return "", time.Time{}, err
	}

	return sessionToken, expiresAt, nil
}

func getIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func pgText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{
			Valid: false,
		}
	}

	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func pgTimestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}
