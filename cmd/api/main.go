package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/mikail-tommard/task-flow/internal/adapters/httpapi"
	"github.com/mikail-tommard/task-flow/internal/adapters/repository"
	"github.com/mikail-tommard/task-flow/internal/adapters/security"
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
	repoAuth := repository.NewUserRepo(db)

	hash := security.NewBcryptHasher(10)

	svc := usecase.NewService(repo)
	svcAuth := usecase.NewAuthService(repoAuth, hash)
	
	mux := httpapi.New(svc, svcAuth)

	handler := httpapi.Chain(
		mux.Routes(),
		httpapi.Recover,
		httpapi.RequestID,
		httpapi.Logging,
	)

	srv := http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		fmt.Println("Starting server on", srv.Addr)
		srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error: ", err)
	}

	_ = db.Close()
}
