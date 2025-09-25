package db_test

import (
	"github.com/gragorther/epigo/database/db"
	"github.com/guregu/null/v6"
)

func (s *Suite) TestCreateGroupsForUser() {
	table := map[string]struct {
		Input []db.CreateGroupsForUser
	}{
		"valid": {
			Input: []db.CreateGroupsForUser{
				{Name: "testname", Description: null.NewString("testdesc", true)},
				{Name: "testname2", Description: null.StringFrom("testdesc")},
			},
		},
	}
	for name, test := range table {
		s.Run(name, func() {
			userID, err := s.Repo.CreateUserReturningID(s.Ctx, db.CreateUserInput{
				Username: "tesetusername",
				Email:    "testemail",
			})
			s.Require().NoError(err, "creating test user shouldn't fail")
			s.Require().NoError(s.Repo.CreateGroupsForUser(s.Ctx, userID, test.Input))

			got, err := s.Repo.GroupsByUserID(s.Ctx, userID)
			s.Require().NoError(err)
			if len(got) != len(test.Input) {
				s.FailNowf("length of got doesn't match input groups", "got: %v, want: %v", len(got), len(test.Input))
				return
			}
			for i, group := range test.Input {
				got := got[i]
				s.Equal(group.Name, got.Name, "group names should match")
				s.Equal(group.Description, got.Description, "descriptions should match")
			}
		})
	}
}

func (s *Suite) TestCanUserEditGroup() {
	s.Run("not authorized", func() {
		userID, err := s.Repo.CreateUserReturningID(s.Ctx, db.CreateUserInput{
			Username:     "testusername",
			Email:        "testemail@gioogle.com",
			PasswordHash: "testhass",
		})
		s.Require().NoError(err, "creating test user shouldn't fail")

		groupID, err := s.Repo.CreateGroupReturningID(s.Ctx, db.CreateGroup{
			UserID: userID,
			Name:   "testname",
		})
		s.Require().NoError(err, "creating test group shouldn't fail")

		authorized, err := s.Repo.CanUserEditGroup(s.Ctx, userID, groupID, nil)
		s.Require().NoError(err, "checking if user can edit group shouldn't fail")
		s.True(authorized, "user should be authorized because they own the group without lastmessages")
	})
}

func (s *Suite) TestGroupByID() {
	table := map[string]db.CreateGroup{
		"valid": {
			Name:        "testname",
			Description: null.StringFrom("test desc"),
		},
	}

	for name, test := range table {
		s.Run(name, func() {
			userID, err := s.Repo.CreateUserReturningID(s.Ctx, db.CreateUserInput{
				Username: "testname",
				Email:    "testemail@google.com",
			})
			s.Require().NoError(err)
			test.UserID = userID
			id, err := s.Repo.CreateGroupReturningID(s.Ctx, test)
			s.Require().NoError(err)

			got, err := s.Repo.GroupByID(s.Ctx, id)
			s.Require().NoError(err)

			s.Equal(test.UserID, got.UserID, "userIDs should match")
			s.Equal(test.Name, got.Name)
			s.Equal(test.Description.String, got.Description.String, "descriptions should match")
			s.Equal(test.LastMessageIDs, got.LastMessageIDs)
		})
	}

	s.Run("last messages", func() {
	})
}
