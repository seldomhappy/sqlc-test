package author

import "context"

// Repository defines the interface for author data access.
type Repository interface {
	GetAuthor(ctx context.Context, id int64) (*Author, error)
	ListAuthors(ctx context.Context) ([]*Author, error)
	CreateAuthor(ctx context.Context, params CreateAuthorParams) (*Author, error)
	UpdateAuthor(ctx context.Context, params UpdateAuthorParams) error
	DeleteAuthor(ctx context.Context, id int64) error
}
