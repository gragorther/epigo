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
