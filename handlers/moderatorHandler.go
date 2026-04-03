package handlers

import (
	"auth-service/mail"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/utils"
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) CreateModerator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	password := utils.GeneratePassword()
	req.Password = password

	if err := repository.CreateUser(h.db, req); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	go func() {
		err := mail.SendTemporaryPasswordEmail(req.Email, password)
		if err != nil {
			log.Println("failed to send temporary code")
		}
	}()

	w.WriteHeader(http.StatusCreated)
}
