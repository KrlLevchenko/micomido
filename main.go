package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/KrlLevchenko/micomido/internal/api"
	"github.com/KrlLevchenko/micomido/internal/repository"
	"github.com/KrlLevchenko/micomido/internal/s3"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is not set")
	}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

	s3Client, err := s3.NewClient(context.Background(), s3Bucket)
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	}
	log.Println("Successfully connected to S3")

	mealRepo := repository.NewMysqlMealRepository(db)
	app := &api.API{
		MealRepo: mealRepo,
		S3Client: s3Client,
	}
	router := api.NewRouter(app)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
