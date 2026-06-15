package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Iamfarhan-cs/crud-app/internal/models"
	"github.com/Iamfarhan-cs/crud-app/internal/storage"
)

type Router struct {
	store *storage.MemoryUserStore
}

func NewRouter(store *storage.MemoryUserStore) http.Handler {
	router := &Router{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", router.health)
	mux.HandleFunc("GET /users", router.listUsers)
	mux.HandleFunc("POST /users", router.createUser)
	mux.HandleFunc("GET /users/{id}", router.getUser)
	mux.HandleFunc("PUT /users/{id}", router.updateUser)
	mux.HandleFunc("DELETE /users/{id}", router.deleteUser)

	return mux
}

func (r *Router) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (r *Router) listUsers(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, r.store.List())
}

func (r *Router) createUser(w http.ResponseWriter, req *http.Request) {
	var input models.UserInput
	if err := readJSON(req, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := r.store.Create(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (r *Router) getUser(w http.ResponseWriter, req *http.Request) {
	user, ok := r.store.Get(req.PathValue("id"))
	if !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (r *Router) updateUser(w http.ResponseWriter, req *http.Request) {
	var input models.UserInput
	if err := readJSON(req, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := r.store.Update(req.PathValue("id"), input)
	if errors.Is(err, storage.ErrUserNotFound) {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (r *Router) deleteUser(w http.ResponseWriter, req *http.Request) {
	if ok := r.store.Delete(req.PathValue("id")); !ok {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func readJSON(req *http.Request, target any) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(message)})
}
