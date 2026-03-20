package httpapi

import (
	"net/http"

	"github.com/mikail-tommard/task-flow/internal/adapters/token"
	"github.com/mikail-tommard/task-flow/internal/usecase"
)

type API struct {
	svc     *usecase.Service
	svcAuth *usecase.AuthService
	jwt     *token.Service
}

func New(svc *usecase.Service, svcAuth *usecase.AuthService, jwt *token.Service) *API {
	return &API{
		svc:     svc,
		svcAuth: svcAuth,
		jwt:     jwt,
	}
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.health)

	mux.Handle("POST /tasks", Chain(http.HandlerFunc(a.createTask),
		Auth(a.jwt),
	))
	mux.HandleFunc("GET /task/{id}", a.getTask)
	mux.HandleFunc("GET /tasks/{userId}", a.listByUser)
	mux.HandleFunc("PATH /task/{id}", a.updateTask)

	mux.HandleFunc("POST /auth/register", a.registerUser)

	return mux
}
