package gormdb_test

import (
	"github.com/gragorther/epigo/models"
	"github.com/samber/lo"
)

func (s *DBTestSuite) TestCheckIfGroupExistsByID() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	s.Run("group doesn't exist", func() {
		exists, err := s.repo.CheckIfGroupExistsByID(s.ctx, 2)
		s.Require().NoError(err)
		s.False(exists, "user shouldn't exist")
	})
	s.Run("group exists", func() {
		group := models.Group{Name: "test group", UserID: userID}
		err := s.repo.CreateGroup(&group)
		s.Require().NoError(err)
		exists, err := s.repo.CheckIfGroupExistsByID(s.ctx, group.ID)
		s.Require().NoError(err)
		s.True(exists, "group should exist")
	})
}

func (s *DBTestSuite) TestDeleteGroupByID() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	group := models.Group{
		Name: "test name", Description: lo.ToPtr("test desc"), Recipients: []models.Recipient{
			{APIRecipient: models.APIRecipient{
				Email: "test@email.com",
			}},
		}, UserID: userID}
	s.Require().NoError(s.repo.CreateGroup(&group))
	s.Require().NoError(s.repo.DeleteGroupByID(s.ctx, group.ID))
	exists, err := s.repo.CheckIfGroupExistsByID(s.ctx, group.ID)
	s.Require().NoError(err)
	s.False(exists, "group should not exist because it was deleted")
	recipientExists, err := s.repo.CheckIfRecipientExistsByID(s.ctx, group.Recipients[0].ID)
	s.Require().NoError(err)
	s.False(recipientExists, "recipient should not exist because their group was deleted")
}

func assertRecipientArrayEquality(s *DBTestSuite, expected, got []models.Recipient) {
	for i := range expected {
		expected := expected[i]
		got := got[i]
		s.Equal(expected.Email, got.Email, "recipient emails should match")
		s.Equal(expected.GroupID, got.GroupID, "group IDs should match")
	}
}

func assertGroupArrayEquality(s *DBTestSuite, expected, got []models.Group) {

	for i := range expected {
		expected := expected[i]
		got := got[i]
		s.Equal(expected.Name, got.Name, "names should match")
		s.Equal(expected.Description, got.Description, "descriptions should match")
		s.Equal(expected.UserID, got.UserID, "userIDs should match")
		assertRecipientArrayEquality(s, expected.Recipients, got.Recipients)
	}
}

func (s *DBTestSuite) TestFindGroupsAndRecipientsByUserID() {
	table := map[string]struct {
		Groups []models.Group
	}{
		"group exists and has recipients": {
			Groups: []models.Group{
				{Name: "test name", Description: lo.ToPtr("test description"), Recipients: []models.Recipient{{APIRecipient: models.APIRecipient{Email: "testemail@email.com"}}, {APIRecipient: models.APIRecipient{Email: "testemail@email2.com"}}}},
				{Name: "test name", Description: lo.ToPtr("test description"), Recipients: []models.Recipient{{APIRecipient: models.APIRecipient{Email: "testemail@email.com"}}, {APIRecipient: models.APIRecipient{Email: "testemail@email2.com"}}}},
			},
		},
	}
	for name, test := range table {
		s.Run(name, func() {
			user := models.User{
				Email: "testemail@emails.com", Username: "test username",
			}
			s.Require().NoError(s.repo.CreateUser(&user))

			for i := range test.Groups {
				test.Groups[i].UserID = user.ID
			}
			s.Require().NoError(s.repo.CreateGroups(s.ctx, &test.Groups))

			got, err := s.repo.FindGroupsAndRecipientsByUserID(s.ctx, user.ID)
			s.Require().NoError(err)
			assertGroupArrayEquality(s, test.Groups, got)

		})
	}
}
