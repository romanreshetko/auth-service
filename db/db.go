package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func New(config Config) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
func ConnectWithRetry(cfg Config) (*sql.DB, error) {
	var database *sql.DB
	var err error

	for i := 1; i <= 10; i++ {
		database, err = New(cfg)
		if err == nil {
			log.Println("connected to database")
			return database, nil
		}

		log.Printf("db not ready, retry %d/10\n", i)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
