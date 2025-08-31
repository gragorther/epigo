package gormdb_test

import (
	"context"

	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/reflectutil"
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
	s.Require().Equal(len(expected), len(got), "lengths of recipient slices should match")

	for i := range expected {
		expected := expected[i]
		got := got[i]
		s.Equal(expected.Email, got.Email, "recipient emails should match")
		if expected.GroupID != 0 {
			s.Equal(expected.GroupID, got.GroupID, "group IDs should match")
		}
	}
}
func assertGroupEquality(s *DBTestSuite, expected, got models.Group) {
	s.Equal(expected.Name, got.Name, "names should match")
	s.Equal(expected.Description, got.Description, "descriptions should match")
	if expected.UserID != 0 {
		s.Equal(expected.UserID, got.UserID, "userIDs should match")
	}
	assertRecipientArrayEquality(s, expected.Recipients, got.Recipients)
	for j, expected := range expected.LastMessages {
		assertLastMessageEquality(s, expected, got.LastMessages[j])
	}
}

func assertGroupArrayEquality(s *DBTestSuite, expected, got []models.Group) {

	for i := range expected {
		expected := expected[i]
		got := got[i]
		assertGroupEquality(s, expected, got)
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
			s.Require().NoError(s.repo.CreateUser(s.ctx, &user))

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

func (s *DBTestSuite) TestCreateGroup() {

	table := map[string]struct {
		Group models.Group
	}{
		"group without lastMessages or recipients": {
			Group: models.Group{
				Name: "test group",
			},
		},
		"group with lastMessages": {
			Group: models.Group{
				Name: "test group with last messages", LastMessages: []models.LastMessage{
					{Title: "test lastmesssage title"},
				},
			},
		},
		"group with recipients and last messages": {
			Group: models.Group{
				Name: "test group name", LastMessages: []models.LastMessage{
					{Title: "test title"},
				},
				Recipients: []models.Recipient{
					{APIRecipient: models.APIRecipient{Email: "test@email.com"}},
				},
			},
		},
	}
	for name, test := range table {
		s.Run(name, func() {
			// reset the db so we don't get username collisions
			s.BeforeTest("DBTestSuite", name)
			user := models.User{
				Email:    "testemailasdf",
				Username: "testusernameasdf",
			}
			s.Require().NoError(s.repo.CreateUser(s.ctx, &user), "creating user shouldn't fail")
			test.Group.UserID = user.ID
			for i := range test.Group.LastMessages {
				test.Group.LastMessages[i].UserID = user.ID
			}

			s.Require().NoError(s.repo.CreateGroup(&test.Group), "creating group shouldn't fail")
			s.NotZero(test.Group.ID, "group ID shouldn't be zero because it should be created when creating group")
			got, err := s.repo.FindGroupsAndLastMessagesAndRecipientsByUserID(s.ctx, user.ID)
			s.Require().NoError(err, "finding user groups shouldn't fail")
			s.Require().Len(got, 1, "the user should have 1 group")
			assertGroupEquality(s, test.Group, got[0])
		})

	}
}

// here, many similar functions are tested
func (s *DBTestSuite) TestFindGroups() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)

	type WantAndFunc struct {
		Func func(ctx context.Context, userID uint) ([]models.Group, error)
		Want []models.Group
	}
	table := struct {
		WantAndFunc []WantAndFunc
		// the groups to be created
		Groups []models.Group
	}{Groups: []models.Group{
		{Name: "testname", UserID: userID, Recipients: []models.Recipient{
			{APIRecipient: models.APIRecipient{Email: "testemail"}},
		}, LastMessages: []models.LastMessage{{Title: "testtitle", UserID: userID, Content: lo.ToPtr("test title")}}},
	},
		WantAndFunc: []WantAndFunc{
			{
				Func: s.repo.FindGroupsAndLastMessagesAndRecipientsByUserID,
				Want: []models.Group{
					{Name: "testname", UserID: userID, Recipients: []models.Recipient{
						{APIRecipient: models.APIRecipient{Email: "testemail"}},
					}, LastMessages: []models.LastMessage{{Title: "testtitle", UserID: userID, Content: lo.ToPtr("test title")}}},
				},
			},
			{
				Func: s.repo.FindGroupsAndLastMessagesByUserID, Want: []models.Group{
					{Name: "testname", UserID: userID, LastMessages: []models.LastMessage{{Title: "testtitle", UserID: userID, Content: lo.ToPtr("test title")}}},
				},
			},
			{
				Func: s.repo.FindGroupsAndRecipientsByUserID, Want: []models.Group{
					{Name: "testname", UserID: userID, Recipients: []models.Recipient{
						{APIRecipient: models.APIRecipient{Email: "testemail"}},
					}},
				},
			},
		},
	}

	s.Require().NoError(s.repo.CreateGroups(s.ctx, &table.Groups))

	for _, want := range table.WantAndFunc {
		s.Run(reflectutil.GetFunctionName(want.Func), func() {
			got, err := want.Func(s.ctx, userID)
			s.Require().NoError(err)
			assertGroupArrayEquality(s, want.Want, got)
		})

	}

}

func (s *DBTestSuite) TestCreateGroups() {

	table := map[string]struct {
		Groups []models.Group
	}{
		"contains last messages": {
			Groups: []models.Group{
				{Name: "testname", LastMessages: []models.LastMessage{
					{Title: "test title"},
				}},
			},
		},
		"contains recipients": {
			Groups: []models.Group{
				{Name: "test name", Recipients: []models.Recipient{
					{APIRecipient: models.APIRecipient{Email: "testemail"}},
				}},
			},
		},
	}
	for name, test := range table {
		s.Run(name, func() {
			s.BeforeTest("DBTestSuite", name)
			userID, err := createTestUser(s.ctx, s.db)
			s.Require().NoError(err)
			for i := range test.Groups {
				test.Groups[i].UserID = userID
				for j := range test.Groups[i].LastMessages {
					test.Groups[i].LastMessages[j].UserID = userID
				}
			}

			s.Require().NoError(s.repo.CreateGroups(s.ctx, &test.Groups))
			got, err := s.repo.FindGroupsAndLastMessagesAndRecipientsByUserID(s.ctx, userID)
			s.Require().NoError(err)
			assertGroupArrayEquality(s, test.Groups, got)
		})
	}
}

func (s *DBTestSuite) TestUpdateGroup() {
	userID, err := createTestUser(s.ctx, s.db)
	s.Require().NoError(err)
	group := models.Group{
		Name:        "testname",
		UserID:      userID,
		Description: lo.ToPtr("test desc"),
		Recipients: []models.Recipient{
			{APIRecipient: models.APIRecipient{Email: "testemail@emal.com"}},
			{APIRecipient: models.APIRecipient{Email: "testemial@email.com"}},
		},
		LastMessages: []models.LastMessage{
			{Title: "test title", Content: lo.ToPtr("test content"), UserID: userID},
		},
	}
	s.Require().NoError(s.repo.CreateGroup(&group))

	updatedGroup := models.Group{
		Name: "testname 2", Description: lo.ToPtr("new desc"), Recipients: []models.Recipient{
			{APIRecipient: models.APIRecipient{Email: "newrecipient"}},
		}, ID: group.ID,
		LastMessages: []models.LastMessage{
			{Title: "new title", Content: lo.ToPtr("new content"), UserID: userID},
		},
	}
	s.Require().NoError(s.repo.UpdateGroup(s.ctx, updatedGroup))
	got, err := s.repo.FindGroupsAndLastMessagesAndRecipientsByUserID(s.ctx, userID)
	s.Require().NoError(err)
	assertGroupEquality(s, updatedGroup, got[0])
}
