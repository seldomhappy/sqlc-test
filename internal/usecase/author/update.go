package author

import (
	"context"

	"github.com/seldomhappy/sqlc-test/internal/domain/author"
)

// UpdateAuthorUseCase updates an existing author.
type UpdateAuthorUseCase struct {
	repo author.Repository
}

// NewUpdateAuthorUseCase creates a new UpdateAuthorUseCase.
func NewUpdateAuthorUseCase(repo author.Repository) *UpdateAuthorUseCase {
	return &UpdateAuthorUseCase{repo: repo}
}

// Execute updates an author.
func (u *UpdateAuthorUseCase) Execute(ctx context.Context, params author.UpdateAuthorParams) error {
	return u.repo.UpdateAuthor(ctx, params)
}
