package author

import (
	"context"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

// CreateAuthorUseCase creates a new author.
type CreateAuthorUseCase struct {
	repo author.Repository
}

// NewCreateAuthorUseCase creates a new CreateAuthorUseCase.
func NewCreateAuthorUseCase(repo author.Repository) *CreateAuthorUseCase {
	return &CreateAuthorUseCase{repo: repo}
}

// Execute creates a new author.
func (u *CreateAuthorUseCase) Execute(ctx context.Context, params author.CreateAuthorParams) (*author.Author, error) {
	return u.repo.CreateAuthor(ctx, params)
}
