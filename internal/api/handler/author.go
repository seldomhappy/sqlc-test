package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
	usecase "github.com/seldomhappy/sqlc-test/internal/usecase/author"
)

// AuthorHandler handles HTTP requests for author operations.
type AuthorHandler struct {
	listUC   *usecase.ListAuthorsUseCase
	getUC    *usecase.GetAuthorUseCase
	createUC *usecase.CreateAuthorUseCase
	updateUC *usecase.UpdateAuthorUseCase
	deleteUC *usecase.DeleteAuthorUseCase
}

// NewAuthorHandler creates a new AuthorHandler.
func NewAuthorHandler(
	listUC *usecase.ListAuthorsUseCase,
	getUC *usecase.GetAuthorUseCase,
	createUC *usecase.CreateAuthorUseCase,
	updateUC *usecase.UpdateAuthorUseCase,
	deleteUC *usecase.DeleteAuthorUseCase,
) *AuthorHandler {
	return &AuthorHandler{
		listUC:   listUC,
		getUC:    getUC,
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

// ListAuthors handles GET /authors.
func (h *AuthorHandler) ListAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authors, err := h.listUC.Execute(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list authors"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authors)
}

// GetAuthor handles GET /authors/{id}.
func (h *AuthorHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid author id"}`, http.StatusBadRequest)
		return
	}

	author, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"author not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(author)
}

// CreateAuthor handles POST /authors.
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params author.CreateAuthorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	created, err := h.createUC.Execute(r.Context(), params)
	if err != nil {
		http.Error(w, `{"error":"failed to create author"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// UpdateAuthor handles PUT /authors/{id}.
func (h *AuthorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid author id"}`, http.StatusBadRequest)
		return
	}

	var params author.UpdateAuthorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}
	params.ID = id

	if err := h.updateUC.Execute(r.Context(), params); err != nil {
		http.Error(w, `{"error":"failed to update author"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAuthor handles DELETE /authors/{id}.
func (h *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid author id"}`, http.StatusBadRequest)
		return
	}

	if err := h.deleteUC.Execute(r.Context(), id); err != nil {
		http.Error(w, `{"error":"failed to delete author"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers all author routes.
func (h *AuthorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /authors", h.ListAuthors)
	mux.HandleFunc("GET /authors/{id}", h.GetAuthor)
	mux.HandleFunc("POST /authors", h.CreateAuthor)
	mux.HandleFunc("PUT /authors/{id}", h.UpdateAuthor)
	mux.HandleFunc("DELETE /authors/{id}", h.DeleteAuthor)
}
