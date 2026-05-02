package http

import (
	"encoding/json"
	"net/http"

	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/domain"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/usecase"
)

type Handler struct {
	uc *usecase.AuthUseCase
}

func NewHandler(uc *usecase.AuthUseCase) *Handler {
	return &Handler{uc: uc}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "username and password required"})
		return
	}

	host := r.Host
	result, err := h.uc.Login(r.Context(), domain.Credentials{
		Username: req.Username,
		Password: req.Password,
	}, host)

	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		case domain.ErrEmptyHost:
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: result.Token})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
