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
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "multipart too large", http.StatusBadRequest)
		return
	}

	request := r.FormValue("request")
	if request == "" {
		http.Error(w, "missing review", http.StatusBadRequest)
		return
	}

	var req models.RegisterRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Photo != nil {
		photo, err := utils.SavePhoto(r, *(req.Photo))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		*(req.Photo) = photo
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

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
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
			return
		}
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	if user.Role != "user" {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func (h *Handler) GetUserByToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	user, err := repository.GetUser(h.db, claims.UserID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
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

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "multipart too large", http.StatusBadRequest)
		return
	}

	request := r.FormValue("request")
	if request == "" {
		http.Error(w, "missing request", http.StatusBadRequest)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Photo != nil {
		photo, err := utils.SavePhoto(r, *(req.Photo))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		*(req.Photo) = photo
	}

	err := repository.UpdateUser(h.db, claims.UserID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateUserPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "service" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect user_id", http.StatusBadRequest)
		return
	}
	points, err := strconv.ParseInt(r.URL.Query().Get("points"), 10, 32)
	if err != nil {
		http.Error(w, "incorrect points value", http.StatusBadRequest)
		return
	}

	err = repository.UpdateUserPoints(h.db, userID, points)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating points", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
