package author

import (
	"context"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

// ListAuthorsUseCase retrieves all authors.
type ListAuthorsUseCase struct {
	repo author.Repository
}

// NewListAuthorsUseCase creates a new ListAuthorsUseCase.
func NewListAuthorsUseCase(repo author.Repository) *ListAuthorsUseCase {
	return &ListAuthorsUseCase{repo: repo}
}

// Execute retrieves all authors.
func (u *ListAuthorsUseCase) Execute(ctx context.Context) ([]*author.Author, error) {
	return u.repo.ListAuthors(ctx)
}
