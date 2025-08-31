package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/mock"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUsername = "username"
const testPass = "testpass"

func TestLoginUser(t *testing.T) {

	hash := func(pass string) string {
		require := require.New(t)
		hash, err := createHash(pass, argon2id.DefaultParams)
		require.NoError(err)
		return hash
	}
	type want struct {
		Status int
		Token  bool // whether the user gets a token
	}
	testUser := models.User{Username: testUsername, PasswordHash: hash(testPass)}
	table := []struct {
		Name         string
		Input        handlers.LoginInput
		UserToCreate models.User
		Want         want
	}{
		{Name: "valid input", Input: handlers.LoginInput{Username: testUsername, Password: testPass}, UserToCreate: testUser, Want: want{Status: http.StatusOK, Token: true}},
		{Name: "no password", Input: handlers.LoginInput{Username: testUser.Username}, UserToCreate: testUser, Want: want{Status: http.StatusUnprocessableEntity}},
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			r := gin.Default()
			mock := mock.NewMockDB()
			_ = mock.CreateUser(t.Context(), &test.UserToCreate)
			r.POST("/", handlers.LoginUser(mock, comparePasswordAndHash, JWT_SECRET, testIssuer, []string{testIssuer}))

			input, err := sonic.MarshalString(test.Input)
			require.NoError(err)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(input))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(test.Want.Status, w.Code, "http status code should indicate success")
			if test.Want.Token {
				responseBytes, err := io.ReadAll(w.Body)
				require.NoError(err, "reading http response shouldn't fail")
				var response handlers.LoginResponse
				require.NoError(sonic.Unmarshal(responseBytes, &response))
				userID, err := tokens.ParseUserAuth(JWT_SECRET, response.Token, testIssuer, []string{testIssuer})
				require.NoError(err)
				assert.Equal(test.UserToCreate.ID, userID, "userIDs should match")
			}
		})
	}
}
