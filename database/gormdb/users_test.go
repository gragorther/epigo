package gormdb_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/database/initializers"
	"github.com/gragorther/epigo/database/testhelpers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repo        *gormdb.GormDB
	ctx         context.Context
	db          *sql.DB
}

func (suite *DBTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	//here we connect and then close the DB, just to run the migrations. The connection
	// *must* be closed, otherwise the snapshot fails because there can't be any active connections.
	conn, err := initializers.ConnectDB(suite.ctx, pgContainer.ConnectionString)
	suite.Require().NoError(err)
	db, err := conn.DB()
	suite.Require().NoError(err)
	db.Close()

	suite.Require().NoError(suite.pgContainer.Snapshot(suite.ctx))

	conn, err = initializers.ConnectDB(suite.ctx, pgContainer.ConnectionString)
	suite.db, err = conn.DB()
	suite.Require().NoError(err)
	// check Migrator for the table

	repo := gormdb.NewGormDB(conn)
	suite.repo = repo
}
func (suite *DBTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
	if err := suite.db.Close(); err != nil {
		log.Fatalf("failed to close sql db: %v", err)
	}
}

func (suite *DBTestSuite) AfterTest(suiteName, testName string) {
	// close existing sql.DB so the restore can succeed cleanly
	if suite.db != nil {
		_ = suite.db.Close()
	}

	err := suite.pgContainer.Restore(suite.ctx)
	suite.Require().NoError(err)

	// reconnect a fresh DB & repo because restore kills previous connections
	conn, err := initializers.ConnectDB(suite.ctx, suite.pgContainer.ConnectionString)
	suite.Require().NoError(err)
	suite.db, err = conn.DB()
	suite.Require().NoError(err)

	suite.repo = gormdb.NewGormDB(conn)
}

func (suite *DBTestSuite) TestCreateProfile() {
	suite.db.ExecContext(suite.ctx, "INSERT INTO users (username) VALUES ('femboy')")

	name := "uwu"
	newProfile := models.Profile{
		Name:   &name,
		UserID: 1,
	}

	suite.repo.CreateProfile(suite.ctx, &newProfile)

	var got models.Profile
	suite.db.QueryRow("SELECT name FROM profiles WHERE id = 1").Scan(&got.Name)

	suite.Equal(newProfile.Name, got.Name)
}

func (suite *DBTestSuite) TestUpdateProfile() {
	name := "uwu"
	var userID int64
	suite.db.QueryRowContext(suite.ctx, "INSERT INTO users (username) VALUES ('femboy') RETURNING id").Scan(&userID)

	oldProfile := models.Profile{
		Name:   &name,
		UserID: uint(userID),
	}
	err := suite.repo.CreateProfile(suite.ctx, &oldProfile)
	suite.NoError(err)

	newName := "owo"
	newProfile := models.Profile{
		Name: &newName,
		ID:   1,
	}
	suite.repo.UpdateProfile(suite.ctx, newProfile)
	suite.NoError(err)

	var got models.Profile
	suite.db.QueryRow("SELECT name FROM profiles WHERE id = 1").Scan(&got.Name)

	suite.Equal(*newProfile.Name, *got.Name)
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
