package repository

import (
	"database/sql"
	"errors"
)

func InsertCode(db *sql.DB, email, code string) error {
	_, err := db.Exec(`
		INSERT INTO email_verifications
		(email, code, created_at, expires_at) 
		VALUES ($1, $2, NOW(), NOW() + interval '10 minutes')
		ON CONFLICT (email)
		DO UPDATE SET code = EXCLUDED.code, 
		created_at = EXCLUDED.created_at, expires_at = EXCLUDED.expires_at
`, email, code)
	return err
}

func VerifyCode(db *sql.DB, email, code string) (bool, error) {
	var id int64

	err := db.QueryRow(`
		SELECT id FROM email_verifications 
		WHERE email = $1 AND code = $2
		AND expires_at > NOW()
`, email, code).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ConfirmEmail(db *sql.DB, email string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE users
		SET email_verified = TRUE
		WHERE email = $1
`, email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM email_verifications
		WHERE email = $1
`, email)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CheckResendCode(db *sql.DB, email string) (bool, error) {
	var id int64

	err := db.QueryRow(`
		SELECT id FROM email_verifications
		WHERE email = $1
		AND NOW() - created_at > interval '60 seconds'
`, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
