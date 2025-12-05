package author

import (
	"context"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

// GetAuthorUseCase retrieves a single author by ID.
type GetAuthorUseCase struct {
	repo author.Repository
}

// NewGetAuthorUseCase creates a new GetAuthorUseCase.
func NewGetAuthorUseCase(repo author.Repository) *GetAuthorUseCase {
	return &GetAuthorUseCase{repo: repo}
}

// Execute retrieves a single author by ID.
func (u *GetAuthorUseCase) Execute(ctx context.Context, id int64) (*author.Author, error) {
	return u.repo.GetAuthor(ctx, id)
}
