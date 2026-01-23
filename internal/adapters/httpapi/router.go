package httpapi

import (
	"net/http"

	"github.com/mikail-tommard/task-flow/internal/usecase"
)

type API struct {
	svc *usecase.Service
}

func New(svc *usecase.Service) *API {
	return &API{
		svc: svc,
	}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.health)

	mux.HandleFunc("POST /tasks", a.createTask)
	mux.HandleFunc("GET /task/{id}", a.getTask)

	return mux
}
