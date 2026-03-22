package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mikail-tommard/task-flow/internal/usecase"
)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type taskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	UserID      int    `json:"user_id"`
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	var req createTaskRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	userID, ok := r.Context().Value(ctxUserID).(int)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	in := usecase.Input{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID,
	}
	t, err := a.svc.CreateTask(r.Context(), in)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (a *API) getTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	t, err := a.svc.GetTask(r.Context(), i)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	resp := taskResponse{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Done:        t.Done(),
		UserID:      t.UserID(),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (a *API) listByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	id := r.PathValue("userId")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	tasks, err := a.svc.ListTasks(r.Context(), i)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	var resp []taskResponse

	for _, t := range tasks {
		resp = append(resp, taskResponse{
			ID:          t.ID(),
			Title:       t.Title(),
			Description: t.Description(),
			Done:        t.Done(),
			UserID:      t.UserID(),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func (a *API) updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	var req updateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	t, err := a.svc.UpdateTask(r.Context(), usecase.UpdateTaskInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	resp := taskResponse{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Done:        t.Done(),
		UserID:      t.UserID(),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (a *API) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
	}

	user, err := a.svcAuth.CreateUser(r.Context(), usecase.InputUser{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, registerUserResponse{
		ID:    user.UserID(),
		Email: user.Email(),
	})
}

func (a *API) loginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method not allowed")
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request", "invalid request")
		return
	}

	token, err := a.svcAuth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid credential", "invalid credential")
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token: token,
	})
}
