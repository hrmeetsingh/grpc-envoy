package http_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httphandler "github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/http"
	jwtadapter "github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/jwt"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/adapter/userstore"
	"github.com/harmeetsingh/grpc-envoy/services/auth/internal/usecase"
)

func setup(t *testing.T) *httphandler.Handler {
	t.Helper()
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	store := userstore.NewHardcoded()
	signer := jwtadapter.NewRS256Signer(key, "https://test.local")
	uc := usecase.New(store, signer)
	return httphandler.NewHandler(uc)
}

func TestLogin_MethodNotAllowed(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	req.Host = "acme.example.com"
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := setup(t)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("not json"))
	req.Host = "acme.example.com"
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	h := setup(t)
	body, _ := json.Marshal(map[string]string{"username": "alice", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Host = "acme.example.com"
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_Success_ReturnsToken(t *testing.T) {
	h := setup(t)
	body, _ := json.Marshal(map[string]string{"username": "alice", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Host = "acme.example.com"
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["token"] == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestLogin_EmptyFields_ReturnsBadRequest(t *testing.T) {
	h := setup(t)
	body, _ := json.Marshal(map[string]string{"username": "", "password": ""})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Host = "acme.example.com"
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

