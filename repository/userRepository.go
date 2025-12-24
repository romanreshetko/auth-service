package repository

import (
	"auth-service/handlers"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, user handlers.RegisterRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	_, err := db.Exec("INSERT INTO users (email, nickname, password, photo, city, status, agreement_pd, agreement_ea) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		user.Email, user.Nickname, hash, user.Photo, user.City, user.Status, user.AgreementPD, user.AgreementEA)
	return err
}

func Authenticate(db *sql.DB, email, password string) (string, error) {
	var id, hash string
	err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&id, &hash)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return id, nil
}
