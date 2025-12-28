package main

import (
	"auth-service/auth"
	"auth-service/handlers"
	"auth-service/utils"
	"database/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := utils.LoadPrivateKey("/keys/private.pem")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := utils.LoadPublicKey("/keys/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := auth.AuthMiddleware(publicKey)

	h := handlers.New(db, privateKey)
	mux := http.NewServeMux()
	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/users", h.GetUser)
	mux.Handle("/users/update", authMiddleware(http.HandlerFunc(h.UpdateProfile)))
	log.Println("Auth service started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
