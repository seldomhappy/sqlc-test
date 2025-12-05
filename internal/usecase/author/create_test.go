package author

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

func TestCreateAuthorUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{},
	}
	uc := NewCreateAuthorUseCase(mockRepo)
	params := author.CreateAuthorParams{
		Name: "Charlie",
		Bio:  pgtype.Text{String: "New author", Valid: true},
	}

	// Act
	created, err := uc.Execute(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created == nil {
		t.Fatal("expected created author, got nil")
	}
	if created.Name != "Charlie" {
		t.Errorf("expected name 'Charlie', got '%s'", created.Name)
	}
	if created.ID == 0 {
		t.Errorf("expected non-zero ID, got %d", created.ID)
	}
}

func TestCreateAuthorUseCase_ExecuteEmptyName(t *testing.T) {
	// Arrange
	mockRepo := &mockRepository{
		authors: []*author.Author{},
	}
	uc := NewCreateAuthorUseCase(mockRepo)
	params := author.CreateAuthorParams{
		Name: "",
		Bio:  pgtype.Text{String: "", Valid: false},
	}

	// Act
	created, err := uc.Execute(context.Background(), params)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Repository will accept empty names, but in production you might validate
	if created.Name != "" {
		t.Errorf("expected empty name, got '%s'", created.Name)
	}
}
