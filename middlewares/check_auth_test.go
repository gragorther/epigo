package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var JWT_SECRET []byte = []byte("very sercure")

type Mock struct {
	Users []models.User
}

func (m *Mock) GetUserByID(ctx context.Context, ID uint) (models.User, error) {
	for _, user := range m.Users {
		if user.ID == ID {
			return user, nil
		}
	}
	return models.User{}, nil
}

const testUserID = 1

func TestCheckAuth(t *testing.T) {

	type want struct {
		Status int
		User   uint
	}
	table := []struct {
		Name   string
		Header http.Header
		Want   want
	}{
		{Name: "valid", Want: want{
			User:   testUserID,
			Status: http.StatusOK,
		}},
		{Name: "invalid token", Want: want{
			Status: http.StatusUnauthorized,
		}},
	}

	{
		require := require.New(t)
		token, err := tokens.CreateUserAuth(JWT_SECRET, testUserID)
		require.NoError(err, "creating user auth token shouldn't fail")

		middlewares.SetHttpAuthHeaderToken(&table[0].Header, token)
	}
	{
		table[1].Header = make(http.Header)
		table[1].Header.Set("Authorization", "Bearer asekfjasiefjhasjeflčaisjefčalsdjf")
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			assert := assert.New(t)
			gin.SetMode(gin.TestMode)

			r := gin.New()
			r.Use(middlewares.CheckAuth(JWT_SECRET))

			userIDs := make(chan uint, 1)

			r.GET("/", func(c *gin.Context) {
				userID, err := handlers.GetUserIDFromContext(c)
				if err != nil {
					userIDs <- 0
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				/*
				 sends the user ID we got from the context into the channel,
				 because otherwise the test can't access the userID and check if it's the correct one
				*/
				userIDs <- userID

				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header = test.Header
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if test.Want.User != 0 {
				userID := <-userIDs
				assert.Equal(test.Want.User, userID, "user IDs should match")
			}
			assert.Equal(test.Want.Status, w.Code, "status codes should match")

		})
	}

}
