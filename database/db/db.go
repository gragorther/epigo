package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db *pgxpool.Pool
}

func NewDB(db *pgxpool.Pool) *DB {
	return &DB{
		db: db,
	}
}
