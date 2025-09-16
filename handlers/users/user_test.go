package users_test

import (
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/db"
	ginctx "github.com/gragorther/epigo/handlers/context"
	"github.com/gragorther/epigo/handlers/httptesthelpers"
	"github.com/gragorther/epigo/handlers/users"
	"github.com/guregu/null/v6"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	httptesthelpers.HandlersTestSuite
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

const JwtSecret = "testsecret"

func (s *UserSuite) TestGetData() {
	table := map[string]struct {
		User db.CreateUserInput
	}{
		"normal user": {
			User: db.CreateUserInput{
				Username:     "testuername",
				Email:        "testemail@google.com",
				PasswordHash: "securehash123",
				Name:         null.NewString("name", true),
			},
		},
		"no name": {
			User: db.CreateUserInput{
				Username:     "testuername",
				Email:        "testemail@google.com",
				PasswordHash: "securehash123",
				//	Name:         null.NewString("name", true),
			},
		},
	}

	for name, test := range table {
		s.Run(name, func() {
			gin.SetMode(gin.TestMode)
			userID, err := s.Repo.CreateUserReturningID(s.Ctx, test.User)
			s.Require().NoError(err, "creating user shouldn't fail")

			c, w := httptesthelpers.CreateTestContext()
			ginctx.SetUserID(c, userID)
			users.GetData(s.Repo)(c)
			s.AssertHTTPStatus(c, http.StatusOK, w)
			var response users.GetUserDataOutput
			s.Require().NoError(sonic.Unmarshal(w.Body.Bytes(), &response))
			s.Equal(test.User.Name.String, response.Name)
			s.Equal(test.User.Email, response.Email)
			s.Equal(test.User.Username, response.Username)
			s.T().Log(c.Errors)
		})
	}
}
