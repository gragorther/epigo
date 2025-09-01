package handlers_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/mock"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testIssuer = "https://server.com"

var parseEmailVerificationToken = tokens.ParseEmailVerification(JWT_SECRET, testIssuer, testIssuer)

func TestRegisterUser(t *testing.T) {
	requi := require.New(t)
	createToken := func(email string, audience string, issuer string) string {
		token, err := tokens.CreateEmailVerification(JWT_SECRET, audience, issuer)(email)
		requi.NoError(err)
		return token
	}

	type want struct {
		Status      int
		UserCreated bool
	}
	table := []struct {
		Name  string
		Input handlers.RegistrationInput
		Want  want
	}{
		{Name: "valid input", Want: want{Status: http.StatusCreated, UserCreated: true}, Input: handlers.RegistrationInput{Username: "testusername", Name: lo.ToPtr("testname"), Password: "very secure pass", Token: createToken("user@user.com", testIssuer, testIssuer)}},
		{Name: "missing token", Want: want{Status: http.StatusUnprocessableEntity, UserCreated: false}, Input: handlers.RegistrationInput{Username: "testusernae2", Password: "securePass"}},
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			mock := mock.NewMockDB()
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/", handlers.RegisterUser(mock, createHash, parseEmailVerificationToken))
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			input, err := sonic.Marshal(test.Input)
			require.NoError(err, "marshalling json shouldn't fail")
			req.Body = io.NopCloser(bytes.NewBuffer(input))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(test.Want.Status, w.Code)

			exists, _ := mock.CheckIfUserExistsByUsername(t.Context(), test.Input.Username)
			assert.Equal(test.Want.UserCreated, exists, "whether the user exists should match the requirements")

			if test.Want.UserCreated {
				require.Len(mock.Users, 1, "user array length should be 1 because there was 1 user created")
				user := mock.Users[0]
				assert.Equal(test.Input.Username, user.Username, "usernames should match")
				hash, _ := createHash(test.Input.Password, argon2id.DefaultParams)
				assert.Equal(hash, user.PasswordHash, "the hashed password should match the one in the database")
				email, err := parseEmailVerificationToken(test.Input.Token)
				require.NoError(err, "parsing email verification token shouldn't fail")
				assert.Equal(email, user.Email)
			}

		})
	}

}

type TaskEnqueuer struct {
	GotTask *asynq.Task
}

func (t *TaskEnqueuer) EnqueueTask(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	t.GotTask = task
	return nil, nil
}

func TestVerifyEmail(t *testing.T) {
	taskEnqueuer := TaskEnqueuer{}
	type want struct {
		Status int
	}
	table := []struct {
		Name  string
		Input handlers.EmailVerificationInput
		Want  want
	}{
		{Name: "valid input", Input: handlers.EmailVerificationInput{Email: "test@testing.com"}, Want: want{Status: http.StatusOK}},
	}
	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			mock := mock.NewMockDB()
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/", handlers.VerifyEmail(taskEnqueuer.EnqueueTask, mock))
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			input, err := sonic.Marshal(test.Input)
			require.NoError(err)
			req.Body = io.NopCloser(bytes.NewBuffer(input))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(test.Want.Status, w.Code)

			// unmarshal task payload
			var payload tasks.VerificationEmailPayload

			require.NoError(sonic.Unmarshal(taskEnqueuer.GotTask.Payload(), &payload))

			assert.Equal(test.Input.Email, payload.Email, "input email and payload email should match")

		})
	}
}
