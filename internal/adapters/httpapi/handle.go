package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/mikail-tommard/task-flow/internal/usecase"
)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	in := usecase.Input{
		Title:       req.Title,
		Description: req.Description,
		UserID:      req.UserID,
	}
	t, err := a.svc.CreateTask(r.Context(), in)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, t)
}
