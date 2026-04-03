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

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	correctCode, err := repository.VerifyCode(h.db, req.Email, req.Code)
	if err != nil {
		http.Error(w, "failed to verify code", http.StatusInternalServerError)
		return
	}
	if !correctCode {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]bool{"verified": correctCode})
		if err != nil {
			return
		}
		return
	}

	err = repository.ConfirmEmail(h.db, req.Email)
	if err != nil {
		http.Error(w, "failed to confirm email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]bool{"verified": correctCode})
	if err != nil {
		return
	}

}

func (h *Handler) ResendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ResendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	checkedResend, err := repository.CheckResendCode(h.db, req.Email)
	if err != nil {
		http.Error(w, "failed to check resend code", http.StatusInternalServerError)
		return
	}

	if !checkedResend {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{"status": "too early request"})
		w.WriteHeader(http.StatusTooEarly)
		return
	}

	code := utils.GenerateVerificationCode()
	if err := repository.InsertCode(h.db, req.Email, code); err != nil {
		http.Error(w, "failed to insert code", http.StatusInternalServerError)
		return
	}
	go func() {
		err := mail.SendVerificationEmail(req.Email, code)
		if err != nil {
			log.Println("failed to send verification email")
		}
	}()

	w.WriteHeader(http.StatusOK)
}
