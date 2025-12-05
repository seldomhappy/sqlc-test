package author

import (
	"context"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

// DeleteAuthorUseCase deletes an author.
type DeleteAuthorUseCase struct {
	repo author.Repository
}

// NewDeleteAuthorUseCase creates a new DeleteAuthorUseCase.
func NewDeleteAuthorUseCase(repo author.Repository) *DeleteAuthorUseCase {
	return &DeleteAuthorUseCase{repo: repo}
}

// Execute deletes an author by ID.
func (u *DeleteAuthorUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.DeleteAuthor(ctx, id)
}
