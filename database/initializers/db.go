package initializers

import (
	"context"

	"github.com/gragorther/epigo/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func ConnectDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, dsn)
}

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	goose.SetBaseFS(migrations.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	sqlDB := stdlib.OpenDBFromPool(db)

	if err := goose.UpContext(ctx, sqlDB, "assets"); err != nil {
		if err := sqlDB.Close(); err != nil {
			return err
		}
		return err
	}

	return sqlDB.Close()
}
