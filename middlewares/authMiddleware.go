package middlewares

import (
	"auth-service/models"
	"auth-service/utils"
	"context"
	"crypto/rsa"
	"net/http"
	"strings"
)

func AuthMiddleware(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, err := utils.ParseToken(parts[1], publicKey)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "claims", models.AuthContext{
				UserID: int64(claims["user_id"].(float64)),
				Role:   claims["role"].(string),
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
