package main

import (
	"database/sql"
	"log"
	"rest-api-go-gin/internal/database"
	"rest-api-go-gin/internal/env"

	_ "github.com/joho/godotenv/autoload" // used to load environment vars
	_ "modernc.org/sqlite"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
