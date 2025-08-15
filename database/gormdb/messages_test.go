package gormdb_test

import (
	"github.com/gragorther/epigo/models"
	"github.com/samber/lo"
)

func (s *DBTestSuite) TestCreateLastMessage() {
	table := []struct {
		Name        string
		LastMessage models.LastMessage
	}{
		{Name: "no groups in model", LastMessage: models.LastMessage{Title: "testtitle", Content: lo.ToPtr("testcontent"), UserID: 123}},
		{Name: "group in model", LastMessage: models.LastMessage{Title: "testtitle", Content: lo.ToPtr("content"), Groups: []models.Group{
			{Name: "groupname"},
		}}},
	}
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	for _, test := range table {
		s.Run(test.Name, func() {

			//the following two things make sure we don't violate foreign key constraints because there was no user assigned to the group or lastMessage
			test.LastMessage.UserID = userID

			for i := range test.LastMessage.Groups {
				test.LastMessage.Groups[i].UserID = userID
			}

			s.Require().NoError(s.repo.CreateLastMessage(s.ctx, &test.LastMessage))

			got, err := s.repo.GetLastMessageByID(s.ctx, test.LastMessage.ID)
			s.Require().NoError(err)

			s.Equal(test.LastMessage.Content, got.Content)
			s.Equal(test.LastMessage.Title, got.Title)
			s.Equal(test.LastMessage.UserID, got.UserID)
		})
	}

}
func (s *DBTestSuite) TestFindLastMessagesByUserID() {
	table := []struct {
		Name             string
		UserLastMessages []models.LastMessage
	}{
		{Name: "user has last messages", UserLastMessages: []models.LastMessage{
			{Title: "testtitle", Content: lo.ToPtr("testcontent")},
		}},
		{Name: "user doesn't have any lastMessages", UserLastMessages: nil},
	}

	for _, test := range table {
		s.Run(test.Name, func() {
			user := &models.User{
				Username: "testusername", Email: "test@testemails.com",
			}
			s.repo.CreateUser(user)
			for i := range test.UserLastMessages {
				test.UserLastMessages[i].UserID = user.ID
			}
			if test.UserLastMessages != nil {
				s.Require().NoError(s.repo.CreateLastMessages(s.ctx, &test.UserLastMessages))
			}

			got, err := s.repo.FindLastMessagesByUserID(user.ID)
			s.Require().NoError(err)
			s.Len(got, len(test.UserLastMessages), "the length of the input userLastMessages and the one we got from the database should match")
		})
	}
}
