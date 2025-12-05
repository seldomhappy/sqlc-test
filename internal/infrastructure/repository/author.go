package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/seldomhappy/sqlc-test/internal/domain/author"
	"github.com/seldomhappy/sqlc-test/tutorial"
)

// AuthorRepository implements the author.Repository interface using PostgreSQL.
type AuthorRepository struct {
	queries *tutorial.Queries
}

// New creates a new AuthorRepository.
func New(conn *pgx.Conn) *AuthorRepository {
	return &AuthorRepository{
		queries: tutorial.New(conn),
	}
}

// GetAuthor retrieves a single author by ID.
func (r *AuthorRepository) GetAuthor(ctx context.Context, id int64) (*author.Author, error) {
	a, err := r.queries.GetAuthor(ctx, id)
	if err != nil {
		return nil, err
	}
	return &author.Author{
		ID:   a.ID,
		Name: a.Name,
		Bio:  a.Bio,
	}, nil
}

// ListAuthors retrieves all authors.
func (r *AuthorRepository) ListAuthors(ctx context.Context) ([]*author.Author, error) {
	authors, err := r.queries.ListAuthors(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*author.Author, len(authors))
	for i, a := range authors {
		result[i] = &author.Author{
			ID:   a.ID,
			Name: a.Name,
			Bio:  a.Bio,
		}
	}
	return result, nil
}

// CreateAuthor creates a new author.
func (r *AuthorRepository) CreateAuthor(ctx context.Context, params author.CreateAuthorParams) (*author.Author, error) {
	created, err := r.queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: params.Name,
		Bio:  params.Bio,
	})
	if err != nil {
		return nil, err
	}
	return &author.Author{
		ID:   created.ID,
		Name: created.Name,
		Bio:  created.Bio,
	}, nil
}

// UpdateAuthor updates an existing author.
func (r *AuthorRepository) UpdateAuthor(ctx context.Context, params author.UpdateAuthorParams) error {
	return r.queries.UpdateAuthor(ctx, tutorial.UpdateAuthorParams{
		ID:   params.ID,
		Name: params.Name,
		Bio:  params.Bio,
	})
}

// DeleteAuthor deletes an author by ID.
func (r *AuthorRepository) DeleteAuthor(ctx context.Context, id int64) error {
	return r.queries.DeleteAuthor(ctx, id)
}
