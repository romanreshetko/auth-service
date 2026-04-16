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
	emailVerified := false
	if user.Role == "moderator" {
		emailVerified = true
	}
	_, err := db.Exec("INSERT INTO users (email, nickname, password, user_role, photo, city, status, points, agreement_pd, agreement_ea, email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		user.Email, user.Nickname, hash, user.Role, SafeDeref(user.Photo), SafeDeref(user.City), SafeDeref(user.Status), 10, SafeDeref(user.AgreementPD), SafeDeref(user.AgreementEA), emailVerified)
	return err
}

func Authenticate(db *sql.DB, email, password string) (int64, string, error) {
	var id int64
	var hash, role string
	var emailVerified bool
	err := db.QueryRow("SELECT id, password, user_role, email_verified FROM users WHERE email = $1", email).Scan(&id, &hash, &role, &emailVerified)
	if err != nil {
		return 0, "", err
	}
	if !emailVerified {
		return 0, "", errors.New("email not verified")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return 0, "", errors.New("invalid password")
	}

	return id, role, nil
}

func GetUser(db *sql.DB, id int64) (models.UserInfo, error) {
	var userInfo models.UserInfo
	err := db.QueryRow("SELECT email, nickname, user_role, photo, city, points, status FROM users WHERE id = $1", id).Scan(
		&userInfo.Email,
		&userInfo.Nickname,
		&userInfo.Role,
		&userInfo.Photo,
		&userInfo.City,
		&userInfo.Points,
		&userInfo.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userInfo, errors.New("user not found")
		}
		return userInfo, err
	}
	return userInfo, nil
}

func UpdateUser(db *sql.DB, userId int64, user models.UpdateProfileRequest) error {
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

	if user.Password != nil {
		query += fmt.Sprintf("password = $%d, ", idx)
		hash, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		args = append(args, hash)
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

func UpdateUserPoints(db *sql.DB, userID, pointsAdd int64) error {
	res, err := db.Exec(`
		UPDATE users
		SET points = points + $1
		WHERE id = $2
`, pointsAdd, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func SafeDeref[T any](v *T) any {
	if v == nil {
		return nil
	}
	return *v
}
