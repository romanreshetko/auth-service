package main

import (
	"auth-service/auth"
	DB "auth-service/db"
	"auth-service/handlers"
	"auth-service/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	cnf := DB.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := DB.ConnectWithRetry(cnf)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := utils.LoadPrivateKey("./keys/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := utils.LoadPublicKey("./keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := auth.AuthMiddleware(publicKey)

	h := handlers.New(db, privateKey)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/verify", h.VerifyEmail)
	mux.HandleFunc("/resend", h.ResendEmail)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/user", h.GetUser)
	mux.Handle("/users/update", authMiddleware(http.HandlerFunc(h.UpdateProfile)))
	log.Println("Auth service started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
