package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// PostgresDB wraps a pgx.Conn for database operations.
type PostgresDB struct {
	conn *pgx.Conn
}

// New creates a new PostgresDB instance.
func New(conn *pgx.Conn) *PostgresDB {
	return &PostgresDB{conn: conn}
}

// GetConn returns the underlying pgx.Conn.
func (db *PostgresDB) GetConn() *pgx.Conn {
	return db.conn
}

// Close closes the database connection.
func (db *PostgresDB) Close(ctx context.Context) error {
	return db.conn.Close(ctx)
}
