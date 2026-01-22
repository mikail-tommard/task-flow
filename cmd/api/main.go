package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/mikail-tommard/task-flow/internal/adapters/httpapi"
	"github.com/mikail-tommard/task-flow/internal/adapters/repository"
	"github.com/mikail-tommard/task-flow/internal/config"
	"github.com/mikail-tommard/task-flow/internal/usecase"
)

func main() {
	cfg := config.New()

	dsn := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)

	svc := usecase.NewService(repo)
	mux := httpapi.New(svc)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux.Routes(),
	}

	fmt.Println("Starting server on", srv.Addr)
	srv.ListenAndServe()
}
