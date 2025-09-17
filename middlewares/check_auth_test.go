package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ginctx "github.com/gragorther/epigo/handlers/context"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	JWT_SECRET         = []byte("very sercure")
	parseUserAuthToken = tokens.ParseUserAuth(JWT_SECRET, []string{testIssuer}, testIssuer)
)

const (
	testUserID = 1
	testIssuer = "https://issuer.com"
)

var createUserAuth = tokens.CreateUserAuth(JWT_SECRET, []string{testIssuer}, testIssuer)

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
		token, err := createUserAuth(testUserID)
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
			r.Use(middlewares.CheckAuth(parseUserAuthToken))

			userIDs := make(chan uint, 1)

			r.GET("/", func(c *gin.Context) {
				userID, err := ginctx.GetUserID(c)
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
