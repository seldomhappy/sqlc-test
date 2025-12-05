package author

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

func TestUpdateAuthorUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{
			{
				ID:   1,
				Name: "Alice",
				Bio:  pgtype.Text{String: "Original", Valid: true},
			},
		},
	}
	uc := NewUpdateAuthorUseCase(mockRepo)
	params := author.UpdateAuthorParams{
		ID:   1,
		Name: "Alice Updated",
		Bio:  pgtype.Text{String: "Updated bio", Valid: true},
	}

	// Act
	err := uc.Execute(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.authors[0].Name != "Alice Updated" {
		t.Errorf("expected name 'Alice Updated', got '%s'", mockRepo.authors[0].Name)
	}
	if mockRepo.authors[0].Bio.String != "Updated bio" {
		t.Errorf("expected bio 'Updated bio', got '%s'", mockRepo.authors[0].Bio.String)
	}
}

func TestUpdateAuthorUseCase_ExecuteNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{},
	}
	uc := NewUpdateAuthorUseCase(mockRepo)
	params := author.UpdateAuthorParams{
		ID:   999,
		Name: "Ghost",
		Bio:  pgtype.Text{String: "Not real", Valid: true},
	}

	// Act
	err := uc.Execute(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mock repo silently ignores non-existent authors
	if len(mockRepo.authors) != 0 {
		t.Errorf("expected no authors, got %d", len(mockRepo.authors))
	}
}
