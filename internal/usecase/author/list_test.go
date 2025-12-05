package author

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

type mockRepository struct {
	authors []*author.Author
	err     error
}

func (m *mockRepository) GetAuthor(ctx context.Context, id int64) (*author.Author, error) {
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

func (m *mockRepository) ListAuthors(ctx context.Context) ([]*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.authors, nil
}

func (m *mockRepository) CreateAuthor(ctx context.Context, params author.CreateAuthorParams) (*author.Author, error) {
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

func (m *mockRepository) UpdateAuthor(ctx context.Context, params author.UpdateAuthorParams) error {
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

func (m *mockRepository) DeleteAuthor(ctx context.Context, id int64) error {
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

func TestListAuthorsUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{
			{
				ID:   1,
				Name: "Alice",
				Bio:  pgtype.Text{String: "Author 1", Valid: true},
			},
			{
				ID:   2,
				Name: "Bob",
				Bio:  pgtype.Text{String: "Author 2", Valid: true},
			},
		},
	}
	uc := NewListAuthorsUseCase(mockRepo)

	// Act
	authors, err := uc.Execute(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(authors))
	}
	if authors[0].Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", authors[0].Name)
	}
}

func TestListAuthorsUseCase_ExecuteEmpty(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{},
	}
	uc := NewListAuthorsUseCase(mockRepo)

	// Act
	authors, err := uc.Execute(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(authors) != 0 {
		t.Errorf("expected 0 authors, got %d", len(authors))
	}
}
