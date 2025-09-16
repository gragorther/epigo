package testhelpers

import (
	"context"
	"log"

	"github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/database/initializers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	PgContainer *PostgresContainer
	Repo        *db.DB
	Ctx         context.Context
	DB          *pgxpool.Pool
}

func (suite *DBTestSuite) SetupSuite() {
	suite.Ctx = context.Background()

	pgContainer, err := CreatePostgresContainer(suite.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.PgContainer = pgContainer

	// here we connect and then close the DB, just to run the migrations. The connection
	// *must* be closed, otherwise the snapshot fails because there can't be any active connections.
	conn, err := initializers.ConnectDB(suite.Ctx, pgContainer.ConnectionString)
	suite.Require().NoError(err)
	initializers.Migrate(suite.Ctx, conn)
	conn.Close()
	suite.Require().NoError(suite.PgContainer.Snapshot(suite.Ctx))

	conn, err = initializers.ConnectDB(suite.Ctx, pgContainer.ConnectionString)
	suite.Require().NoError(err)
	suite.DB = conn
	suite.Require().NoError(err)

	repo := db.NewDB(conn)
	suite.Repo = repo
}

func (suite *DBTestSuite) TearDownSuite() {
	if err := suite.PgContainer.Terminate(suite.Ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	suite.DB.Close()
}

func setupTest(suite *DBTestSuite) {
	err := suite.PgContainer.Restore(suite.Ctx)
	suite.Require().NoError(err)

	// reconnect a fresh DB & repo because restore kills previous connections
	conn, err := initializers.ConnectDB(suite.Ctx, suite.PgContainer.ConnectionString)
	suite.Require().NoError(err)
	suite.DB = conn

	suite.Repo = db.NewDB(conn)
}

func (suite *DBTestSuite) BeforeTest(suiteName, testName string) {
	setupTest(suite)
}

func (suite *DBTestSuite) SetupSubTest() {
	setupTest(suite)
}
