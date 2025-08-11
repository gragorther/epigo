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

func TestDB(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
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
	suite.Require().NoError(initializers.Migrate(conn))
	sqldb, err := conn.DB()
	suite.Require().NoError(err)
	sqldb.Close()

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

func setupTest(suite *DBTestSuite) {
	err := suite.pgContainer.Restore(suite.ctx)
	suite.Require().NoError(err)

	// reconnect a fresh DB & repo because restore kills previous connections
	conn, err := initializers.ConnectDB(suite.ctx, suite.pgContainer.ConnectionString)
	suite.Require().NoError(err)
	suite.db, err = conn.DB()
	suite.Require().NoError(err)

	suite.repo = gormdb.NewGormDB(conn)
}

func (suite *DBTestSuite) BeforeTest(suiteName, testName string) {
	setupTest(suite)
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
	userID, err := createTestUser(suite.ctx, suite.db)

	oldProfile := models.Profile{
		Name:   &name,
		UserID: uint(userID),
	}
	err = suite.repo.CreateProfile(suite.ctx, &oldProfile)
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

func createTestUser(ctx context.Context, db *sql.DB) (userID uint, err error) {
	err = db.QueryRowContext(ctx, "INSERT INTO users (username) VALUES ('femboy') RETURNING id").Scan(&userID)
	return
}

func (s *DBTestSuite) TestUpdateUserInterval() {
	userID, err := createTestUser(s.ctx, s.db)
	s.NoError(err)
	userCron := "5 4 * * *"
	s.NoError(s.repo.UpdateUserInterval(userID, userCron))

	var gotCron string
	s.NoError(s.db.QueryRowContext(s.ctx, "SELECT cron FROM users WHERE id = $1", userID).Scan(&gotCron))

	s.Equal(userCron, gotCron)
}

func (s *DBTestSuite) TestCheckIfUserExistsByUsernameAndEmail() {
	s.Run("user exists", func() {
		username := "bruh"
		email := "thing@google.com"
		var userID uint
		err := s.db.QueryRowContext(s.ctx, "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id", username, email).Scan(&userID)
		s.Require().NoError(err)

		exists, err := s.repo.CheckIfUserExistsByUsernameAndEmail(username, email)
		s.Require().NoError(err)

		s.True(exists, "should be true because the user does exist")
	})
	s.Run("user doesn't exist", func() {

		exists, err := s.repo.CheckIfUserExistsByUsernameAndEmail("idontexist", "idontexist@test.com")
		s.Require().NoError(err)

		s.False(exists, "should be true because the user does exist")
	})

}

func (s *DBTestSuite) TestCheckIfUserExistsByUsername() {
	s.Run("user exists", func() {
		username := "bruh"
		var userID uint
		err := s.db.QueryRowContext(s.ctx, "INSERT INTO users (username) VALUES ($1) RETURNING id", username).Scan(&userID)
		s.Require().NoError(err)

		exists, err := s.repo.CheckIfUserExistsByUsername(username)
		s.Require().NoError(err)

		s.True(exists, "should be true because the user does exist")
	})
	s.Run("user doesn't exist", func() {

		exists, err := s.repo.CheckIfUserExistsByUsername("idontexist")
		s.Require().NoError(err)

		s.False(exists, "should be true because the user does exist")
	})
}

func (s *DBTestSuite) TestCreateUser() {
	newUser := &models.User{
		Username: "name",
		Email:    "email@email.com",
	}
	s.Require().NoError(s.repo.CreateUser(newUser))

	var username, email string
	err := s.db.QueryRowContext(s.ctx, "SELECT username, email FROM users WHERE id = $1", newUser.ID).Scan(&username, &email)
	s.Require().NoError(err)
	s.Equal(newUser.Username, username)
	s.Equal(newUser.Email, email)
}

func (s *DBTestSuite) TestGetUserByUsername() {
	username := "ime"
	email := "email@email.com"
	s.db.QueryRowContext(s.ctx, "INSERT INTO users (username, email) VALUES ($1, $2)", username, email)

	user, err := s.repo.GetUserByUsername(username)
	s.Require().NoError(err)

	s.Equal(username, user.Username)
	s.Equal(email, user.Email)
}
