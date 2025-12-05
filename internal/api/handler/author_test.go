package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
	usecase "github.com/seldomhappy/sqlc-test/internal/usecase/author"
)

type mockRepoForHandler struct {
	authors []*author.Author
	err     error
}

func (m *mockRepoForHandler) GetAuthor(ctx context.Context, id int64) (*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, a := range m.authors {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, nil
}

func (m *mockRepoForHandler) ListAuthors(ctx context.Context) ([]*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.authors, nil
}

func (m *mockRepoForHandler) CreateAuthor(ctx context.Context, params author.CreateAuthorParams) (*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	a := &author.Author{
		ID:   1,
		Name: params.Name,
		Bio:  params.Bio,
	}
	m.authors = append(m.authors, a)
	return a, nil
}

func (m *mockRepoForHandler) UpdateAuthor(ctx context.Context, params author.UpdateAuthorParams) error {
	if m.err != nil {
		return m.err
	}
	for i, a := range m.authors {
		if a.ID == params.ID {
			m.authors[i].Name = params.Name
			m.authors[i].Bio = params.Bio
			return nil
		}
	}
	return nil
}

func (m *mockRepoForHandler) DeleteAuthor(ctx context.Context, id int64) error {
	if m.err != nil {
		return m.err
	}
	for i, a := range m.authors {
		if a.ID == id {
			m.authors = append(m.authors[:i], m.authors[i+1:]...)
			return nil
		}
	}
	return nil
}

func setupHandler() *AuthorHandler {
	mockRepo := &mockRepoForHandler{
		authors: []*author.Author{
			{
				ID:   1,
				Name: "Alice",
				Bio:  pgtype.Text{String: "Author 1", Valid: true},
			},
		},
	}

	listUC := usecase.NewListAuthorsUseCase(mockRepo)
	getUC := usecase.NewGetAuthorUseCase(mockRepo)
	createUC := usecase.NewCreateAuthorUseCase(mockRepo)
	updateUC := usecase.NewUpdateAuthorUseCase(mockRepo)
	deleteUC := usecase.NewDeleteAuthorUseCase(mockRepo)

	return NewAuthorHandler(listUC, getUC, createUC, updateUC, deleteUC)
}

func TestListAuthors_Success(t *testing.T) {
	// Arrange
	handler := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/authors", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ListAuthors(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestGetAuthor_Success(t *testing.T) {
	// Arrange
	handler := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/authors/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// Act
	handler.GetAuthor(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var result author.Author
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", result.Name)
	}
}

func TestGetAuthor_InvalidID(t *testing.T) {
	// Arrange
	handler := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/authors/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	// Act
	handler.GetAuthor(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateAuthor_Success(t *testing.T) {
	// Arrange
	handler := setupHandler()
	payload := author.CreateAuthorParams{
		Name: "Charlie",
		Bio:  pgtype.Text{String: "New author", Valid: true},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Act
	handler.CreateAuthor(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
	var result author.Author
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Name != "Charlie" {
		t.Errorf("expected name 'Charlie', got '%s'", result.Name)
	}
}

func TestCreateAuthor_InvalidPayload(t *testing.T) {
	// Arrange
	handler := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/authors", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	// Act
	handler.CreateAuthor(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateAuthor_Success(t *testing.T) {
	// Arrange
	handler := setupHandler()
	payload := author.UpdateAuthorParams{
		ID:   1,
		Name: "Alice Updated",
		Bio:  pgtype.Text{String: "Updated", Valid: true},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/authors/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// Act
	handler.UpdateAuthor(w, req)

	// Assert
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestDeleteAuthor_Success(t *testing.T) {
	// Arrange
	handler := setupHandler()
	req := httptest.NewRequest(http.MethodDelete, "/authors/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// Act
	handler.DeleteAuthor(w, req)

	// Assert
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}
