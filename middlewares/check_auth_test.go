package middlewares_test

import (
	"context"
	"fmt"
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

const JWT_SECRET string = "very sercure"

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

func TestCheckAuth(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		const userID = 1
		assert := assert.New(t)
		require := require.New(t)
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		token, err := tokens.CreateUserAuth(JWT_SECRET, userID)
		require.NoError(err, "creating user auth token shouldn't fail")
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		//c.Header("Authorization", fmt.Sprint("Bearer ", token))
		c.Request.Header.Set("Authorization", fmt.Sprint("Bearer ", token))

		middlewares.CheckAuth(JWT_SECRET)(c)

		assert.Equal(http.StatusOK, w.Code, "http status code should indicate success")
		gotID, err := handlers.GetUserIDFromContext(c)
		require.NoError(err, "getting user from context shouldn't fail")
		assert.Equal(uint(userID), gotID, "userID should match the userID check auth stored into the context")
	})

}
