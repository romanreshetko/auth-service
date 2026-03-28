package handlers

import (
	"auth-service/mail"
	"auth-service/models"
	"auth-service/repository"
	"auth-service/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	userID, role, err := repository.Authenticate(h.db, req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			http.Error(w, "invalid password", http.StatusUnauthorized)
			return
		}
		if err.Error() == "email not verified" {
			http.Error(w, "email not verified", http.StatusUnauthorized)
			return
		}
		http.Error(w, "failed to authenticate", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(userID, role, h.privateKey)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
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
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" && claims.Role != "moderator" && claims.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := repository.UpdateUser(h.db, claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
