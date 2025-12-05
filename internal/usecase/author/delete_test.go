package author

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

func TestDeleteAuthorUseCase_Execute(t *testing.T) {
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
	uc := NewDeleteAuthorUseCase(mockRepo)

	// Act
	err := uc.Execute(context.Background(), 1)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mockRepo.authors) != 1 {
		t.Errorf("expected 1 author, got %d", len(mockRepo.authors))
	}
	if mockRepo.authors[0].ID != 2 {
		t.Errorf("expected remaining author with ID 2, got %d", mockRepo.authors[0].ID)
	}
}

func TestDeleteAuthorUseCase_ExecuteNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{
			{
				ID:   1,
				Name: "Alice",
				Bio:  pgtype.Text{String: "Author 1", Valid: true},
			},
		},
	}
	uc := NewDeleteAuthorUseCase(mockRepo)

	// Act
	err := uc.Execute(context.Background(), 999)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Mock repo silently ignores non-existent authors
	if len(mockRepo.authors) != 1 {
		t.Errorf("expected 1 author, got %d", len(mockRepo.authors))
	}
}
