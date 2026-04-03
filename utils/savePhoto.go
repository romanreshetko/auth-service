package utils

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func SavePhoto(r *http.Request, photo string) (string, error) {
	userID := uuid.New().String()
	baseDir := filepath.Join("uploads", "users", userID)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", errors.New("failed to create directory")
	}

	formKey := "photo_" + photo
	file, header, err := r.FormFile(formKey)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("missing photo")
	}

	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := photo + uuid.New().String() + ext
	fullPath := filepath.Join(baseDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", errors.New("cannot save photo")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", errors.New("cannot write photo")
	}

	return userID + "/" + filename, nil
}
