package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePassword() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
