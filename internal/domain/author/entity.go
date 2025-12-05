package author

import "github.com/jackc/pgx/v5/pgtype"

// Author represents an author in the domain model.
type Author struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}

// CreateAuthorParams holds parameters for creating an author.
type CreateAuthorParams struct {
	Name string
	Bio  pgtype.Text
}

// UpdateAuthorParams holds parameters for updating an author.
type UpdateAuthorParams struct {
	ID   int64
	Name string
	Bio  pgtype.Text
}
