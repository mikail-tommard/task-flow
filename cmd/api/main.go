package main

import (
	"fmt"
	"net/http"

	"github.com/mikail-tommard/task-flow/internal/adapters/httpapi"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpapi.Health)

	srv := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on", srv.Addr)
	srv.ListenAndServe()
}