package gormdb_test

import "github.com/gragorther/epigo/models"

func (s *DBTestSuite) TestCheckUserAuthorizationForGroups() {
	s.Run("user is authorized", func() {
		user := models.User{
			Username: "testusername", Email: "testenail@google.om", Groups: []models.Group{
				{Name: "testgroupname"},
				{Name: "testgroupnme2"},
			},
		}
		s.Require().NoError(s.repo.CreateUser(s.ctx, &user))

		authorized, err := s.repo.CheckUserAuthorizationForGroups(s.ctx, []uint{user.Groups[0].ID, user.Groups[1].ID}, user.ID)
		s.Require().NoError(err)
		s.True(authorized, "user should be authorized because they own the groups")
	})
	s.Run("user is not authorized", func() {
		s.BeforeTest("DBtestSuite", "user is not authorized")

		// the user we're checking authorization for
		user := models.User{Username: "testusername", Email: "testemail", Groups: []models.Group{
			{Name: "testname"},
			{Name: "testname2"},
		}}
		s.Require().NoError(s.repo.CreateUser(s.ctx, &user))

		otherUser := models.User{
			Username: "myname", Email: "email@email.email", Groups: []models.Group{
				{Name: "testgroupname"},
				{Name: "testgroupname2"},
			},
		}
		s.Require().NoError(s.repo.CreateUser(s.ctx, &otherUser))
		authorized, err := s.repo.CheckUserAuthorizationForGroups(s.ctx, []uint{user.Groups[0].ID, user.Groups[1].ID, otherUser.Groups[0].ID, otherUser.Groups[1].ID}, user.ID)
		s.Require().NoError(err)
		s.False(authorized, "user shouldn't be authozired because they don't own some of the groups they're trying to access")
	})
}

func (s *DBTestSuite) TestCheckUserAuthorizationForLastMessage() {
	s.Run("authorized", func() {
		user := models.User{Username: "testname", Email: "testemail", LastMessages: []models.LastMessage{
			{Title: "test title"},
			{Title: "test title 2"},
		}}
		s.Require().NoError(s.repo.CreateUser(s.ctx, &user))

		authorized, err := s.repo.CheckUserAuthorizationForLastMessage(s.ctx, user.LastMessages[0].ID, user.ID)
		s.Require().NoError(err)
		s.True(authorized)
	})
	s.Run("unauthorized", func() {
		s.BeforeTest("DBTestSuite", "unauthorized")
		user := models.User{
			Username: "testusername", Email: "testemail",
		}
		s.Require().NoError(s.repo.CreateUser(s.ctx, &user))

		otherUser := models.User{Username: "username2", Email: "testemail2@email.com", LastMessages: []models.LastMessage{
			{Title: "testtitle212q3"},
		}}

		s.Require().NoError(s.repo.CreateUser(s.ctx, &otherUser))

		authorized, err := s.repo.CheckUserAuthorizationForLastMessage(s.ctx, otherUser.LastMessages[0].ID, user.ID)
		s.Require().NoError(err)
		s.False(authorized, "user shouldn't be authorized because another person owns the last message they are trying to access")

	})
}
