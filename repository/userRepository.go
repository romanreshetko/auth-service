package repository

import (
	"auth-service/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, user models.RegisterRequest) error {
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

func GetUser(db *sql.DB, id string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	err := db.QueryRow("SELECT email, nickname, photo, city, status FROM users WHERE id = $1", id).Scan(
		&userInfo.Email,
		&userInfo.Nickname,
		&userInfo.Photo,
		&userInfo.City,
		&userInfo.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userInfo, errors.New("user not found")
		}
		return userInfo, err
	}
	return userInfo, nil
}

func UpdateUser(db *sql.DB, userId string, user models.UpdateProfileRequest) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	idx := 1

	if user.Nickname != nil {
		query += fmt.Sprintf("nickname = $%d, ", idx)
		args = append(args, *user.Nickname)
		idx++
	}

	if user.Photo != nil {
		query += fmt.Sprintf("photo = $%d, ", idx)
		args = append(args, *user.Photo)
		idx++
	}

	if user.City != nil {
		query += fmt.Sprintf("city = $%d, ", idx)
		args = append(args, *user.City)
		idx++
	}

	if user.Status != nil {
		query += fmt.Sprintf("status = $%d, ", idx)
		args = append(args, *user.Status)
		idx++
	}

	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, userId)

	_, err := db.Exec(query, args...)
	return err
}
