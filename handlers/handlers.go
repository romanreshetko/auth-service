package handlers

import (
	"auth-service/models"
	"auth-service/repository"
	"auth-service/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

type Handler struct {
	db         *sql.DB
	privateKey interface{}
}

func New(db *sql.DB, privateKey interface{}) *Handler {
	return &Handler{db, privateKey}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := repository.CreateUser(h.db, req); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	userID, err := repository.Authenticate(h.db, req.Email, req.Password)
	if err != nil {
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(userID, h.privateKey)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	user, err := repository.GetUser(h.db, userId)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "user not found", http.StatusNotFound)
		}
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := repository.UpdateUser(h.db, userId, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
