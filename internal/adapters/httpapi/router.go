package httpapi

import (
	"net/http"

	"github.com/mikail-tommard/task-flow/internal/usecase"
)

type API struct {
	svc *usecase.Service
	svcAuth *usecase.AuthService
}

func New(svc *usecase.Service, svcAuth *usecase.AuthService) *API {
	return &API{
		svc: svc,
		svcAuth: svcAuth,
	}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.health)

	mux.HandleFunc("POST /tasks", a.createTask)
	mux.HandleFunc("GET /task/{id}", a.getTask)
	mux.HandleFunc("GET /tasks/{userId}", a.listByUser)
	mux.HandleFunc("PATH /task/{id}", a.updateTask)

	mux.HandleFunc("POST /auth/register", a.registerUser)

	return mux
}
