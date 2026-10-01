package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/ARKremlin/go_dnd/internal/pkg/logger"
	"github.com/ARKremlin/go_dnd/internal/transport/rest/middleware"
	"github.com/ARKremlin/go_dnd/internal/usecase"
)

type AuthHandler struct {
	auth *usecase.AuthUseCase
}

func NewAuthHandler(auth *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := h.auth.RegisterDM(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidInput):
			respondError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, usecase.ErrUsernameTaken):
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondInternalError(w, r, err)
		}
		return
	}
	respondJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	tok, err := h.auth.LoginDM(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			respondError(w, http.StatusUnauthorized, "invalid credentials")
		default:
			respondInternalError(w, r, err)
		}
		return
	}
	respondJSON(w, http.StatusOK, loginResponse{Token: tok})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.auth.GetMe(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			respondError(w, http.StatusUnauthorized, "unauthorized")
		default:
			respondInternalError(w, r, err)
		}
		return
	}
	respondJSON(w, http.StatusOK, toUserResponse(u))
}

func toUserResponse(u *domain.User) userResponse {
	resp := userResponse{
		ID:   u.ID.String(),
		Role: string(u.Role),
	}
	if u.Username != nil {
		resp.Username = *u.Username
	}
	resp.CreatedAt = u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z")
	return resp
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, errorResponse{Error: msg})
}

func respondInternalError(w http.ResponseWriter, r *http.Request, err error) {
	logger.FromContext(r.Context()).Error("internal error",
		"err", err,
		"method", r.Method,
		"path", r.URL.Path,
	)
	respondError(w, http.StatusInternalServerError, "internal server error")
}
