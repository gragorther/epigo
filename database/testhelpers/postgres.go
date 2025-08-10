package testhelpers

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	ctr, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithDatabase("epigo"),
		postgres.WithUsername("epigo"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)

	// Run any migrations on the database

	// 2. Create a snapshot of the database to restore later
	// tt.options comes the test case, it can be specified as e.g. `postgres.WithSnapshotName("custom-snapshot")` or omitted, to use default name

	dbURL, err := ctr.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	return &PostgresContainer{
		PostgresContainer: ctr,
		ConnectionString:  dbURL,
	}, nil
}
