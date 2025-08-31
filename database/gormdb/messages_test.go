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
		Want             []models.LastMessage
	}{
		{Name: "user has last messages", UserLastMessages: []models.LastMessage{
			{Title: "testtitle", Content: lo.ToPtr("testcontent")},
		}, Want: []models.LastMessage{{Title: "testtitle", Content: lo.ToPtr("testcontent")}}},
		{Name: "user doesn't have any lastMessages", UserLastMessages: nil, Want: nil},
		{Name: "model has groups", UserLastMessages: []models.LastMessage{
			{Title: "testtitle", Content: lo.ToPtr("testcontent"), Groups: []models.Group{
				{Name: "testgroupanme", Description: lo.ToPtr("testdesc")},
			}},
		}, Want: []models.LastMessage{ // shouldn't return groups
			{Title: "testtitle", Content: lo.ToPtr("testcontent")}}}}
	user := &models.User{
		Username: "testusername", Email: "test@testemails.com",
	}
	s.Require().NoError(s.repo.CreateUser(s.ctx, user))
	for _, test := range table {
		s.Run(test.Name, func() {
			// make sure we don't get fkey errors
			for i := range test.UserLastMessages {
				test.UserLastMessages[i].UserID = user.ID
				for j := range test.UserLastMessages[i].Groups {
					test.UserLastMessages[i].Groups[j].UserID = user.ID
				}
			}
			for i := range test.Want {
				test.Want[i].UserID = user.ID
			}

			if test.UserLastMessages != nil {
				s.Require().NoError(s.repo.CreateLastMessages(s.ctx, &test.UserLastMessages))
			}

			got, err := s.repo.FindLastMessagesByUserID(user.ID)
			s.Require().NoError(err)

			equal := s.Equal
			for i, message := range test.Want {
				gotMessage := got[i]
				equal(message.Title, gotMessage.Title, "titles should match")
				equal(message.Content, gotMessage.Content, "content should match")
				equal(message.UserID, gotMessage.UserID, "user IDs should match")
			}
		})
	}
}

func (s *DBTestSuite) TestCreateLastMessages() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err, "creating test user shouldn't fail")
	table := map[string]struct {
		// the last messages to be created
		LastMessages []models.LastMessage
		Want         []models.LastMessage
	}{
		"contains lastmessages": {
			LastMessages: []models.LastMessage{
				{Title: "testtitle", UserID: userID}, {UserID: userID, Title: "testtitle2"},
			}, Want: []models.LastMessage{
				{Title: "testtitle", UserID: userID}, {UserID: userID, Title: "testtitle2"},
			},
		},
	}

	for name, test := range table {
		s.Run(name, func() {
			lastMessages := test.LastMessages
			s.Require().NoError(s.repo.CreateLastMessages(s.ctx, &lastMessages), "creating last messages shouldn't fail")
			got, err := s.repo.FindLastMessagesByUserID(userID)
			s.Require().NoError(err, "finding last messages by user ID shouldn't fail")

			for i, want := range test.Want {
				got := got[i]
				assertLastMessageEquality(s, want, got)
			}
		})
	}
}

func assertLastMessageEquality(s *DBTestSuite, expected models.LastMessage, actual models.LastMessage) {
	s.Equal(expected.Title, actual.Title, "titles should match")
	s.Equal(expected.Content, actual.Content, "contents should match")

}

func (s *DBTestSuite) TestUpdateLastMessage() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err, "creating test user shouldn't fail")
	table := map[string]struct {
		OldMessage models.LastMessage
		NewMessage models.LastMessage
	}{
		"has last message": {
			OldMessage: models.LastMessage{
				Title: "old testtitle", Content: lo.ToPtr("old content"), UserID: userID,
			},
			NewMessage: models.LastMessage{
				Title: "new title", Content: lo.ToPtr("new content"), UserID: userID,
			},
		},
	}
	for name, test := range table {
		s.Run(name, func() {
			require := s.Require()
			// first, we create the old last message which will then be updated
			require.NoError(s.repo.CreateLastMessage(s.ctx, &test.OldMessage))
			test.NewMessage.ID = test.OldMessage.ID
			s.Require().NoError(s.repo.UpdateLastMessage(s.ctx, test.NewMessage), "updating last message shouldn't fail")

			got, err := s.repo.GetLastMessageByID(s.ctx, test.NewMessage.ID)
			require.NoError(err)

			assertLastMessageEquality(s, test.NewMessage, got)
		})
	}
}

func (s *DBTestSuite) TestDeleteLastMessageByID() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	table := map[string]struct {
		Message models.LastMessage
	}{
		"has last message": {
			Message: models.LastMessage{Title: "test title", Content: lo.ToPtr("test content"), UserID: userID},
		},
	}

	for name, test := range table {
		s.Run(name, func() {
			s.Require().NoError(s.repo.CreateLastMessage(s.ctx, &test.Message))
			s.Require().NoError(s.repo.DeleteLastMessageByID(test.Message.ID))

			exists, err := s.repo.CheckIfLastMessageExistsByID(s.ctx, test.Message.ID)
			s.Require().NoError(err, "checking if last message exists shouldn't fail")
			s.False(exists, "last message shouldn't exist")
		})
	}
}

func (s *DBTestSuite) TestCheckIfLastMessageExists() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)

	s.Run("exists", func() {
		lastMessage := models.LastMessage{
			Title: "test title", UserID: userID, Content: lo.ToPtr("testcontent"),
		}
		s.Require().NoError(s.repo.CreateLastMessage(s.ctx, &lastMessage))

		exists, err := s.repo.CheckIfLastMessageExistsByID(s.ctx, lastMessage.ID)
		s.Require().NoError(err)
		s.True(exists, "last message should exist")
	})
	s.Run("doesn't exist", func() {
		exists, err := s.repo.CheckIfLastMessageExistsByID(s.ctx, 999990)
		s.Require().NoError(err)
		s.False(exists, "last message shouldn't exist")
	})
}

func (s *DBTestSuite) TestGetLastMessageByID() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	s.Run("exists", func() {
		lastMessage := models.LastMessage{
			Title: "testtitle", Content: lo.ToPtr("test content"), UserID: userID,
		}
		s.Require().NoError(s.repo.CreateLastMessage(s.ctx, &lastMessage))

		got, err := s.repo.GetLastMessageByID(s.ctx, lastMessage.ID)
		s.Require().NoError(err)
		assertLastMessageEquality(s, lastMessage, got)
	})

}
