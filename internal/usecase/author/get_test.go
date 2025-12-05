package author

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

func TestGetAuthorUseCase_Execute(t *testing.T) {
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
	uc := NewGetAuthorUseCase(mockRepo)

	// Act
	a, err := uc.Execute(context.Background(), 1)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected author, got nil")
	}
	if a.ID != 1 {
		t.Errorf("expected ID 1, got %d", a.ID)
	}
	if a.Name != "Alice" {
		t.Errorf("expected name 'Alice', got '%s'", a.Name)
	}
}

func TestGetAuthorUseCase_ExecuteNotFound(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{},
	}
	uc := NewGetAuthorUseCase(mockRepo)

	// Act
	a, err := uc.Execute(context.Background(), 999)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != nil {
		t.Errorf("expected nil, got author with ID %d", a.ID)
	}
}
