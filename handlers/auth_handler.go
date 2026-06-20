//Package handlers berfungsi menangani request HTTP (lapisan Controller).
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"Go-golang/auth"
	"Go-golang/models"
	"Go-golang/store"
)

type AuthHandler struct {
	store *store.MemStore
}

func NewAuthHandler(s *store.MemStore) *AuthHandler {
	return &AuthHandler{store: s}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

//POST/register user
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body JSON tidak valid"})
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username & password wajib diisi"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal memproses password"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}
	if err := h.store.Create(user); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "registrasi berhasil",
		"user":    user,
	})
}

//POST/login user
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body JSON tidak valid"})
		return
	}

	user, err := h.store.GetByUsername(req.Username)
	if err != nil {

		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "username atau password salah"})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "username atau password salah"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "gagal membuat token"})
		return
	}

	writeJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

//GET/profile
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Context().Value(userContextKey).(string)

	user, err := h.store.GetByUsername(username)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}