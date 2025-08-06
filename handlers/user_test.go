package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

func assertHTTPStatus(t *testing.T, c *gin.Context, expected int, w *httptest.ResponseRecorder, message string) {
	t.Helper()
	assert := assert.New(t)
	// to make sure the http header is actually written.
	c.Writer.WriteHeaderNow()

	assert.Equal(expected, w.Code, message)
}

func (m *mockDB) CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error) {

	for _, user := range m.Users {
		if user.Username == username && user.Email == email {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockDB) CreateUser(newUser *models.User) error {
	m.Users = append(m.Users, *newUser)
	return nil
}

// a predefined hash so tests are much faster to run because they don't have to *actually* hash the password
const predefinedHash string = "verysecurehash"

// a mock of argon2id.CreateHash
func createHash(password string, params *argon2id.Params) (string, error) {
	return predefinedHash, nil
}

func TestRegisterUser(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		username := "mark"
		name := "Down"
		email := "test@google.com"
		password := "5UP3RS3CR37"
		input, err := sonic.Marshal(handlers.RegistrationInput{
			Username: username,
			Name:     &name,
			Email:    email,
			Password: password,
		})
		if err != nil {
			t.Fatalf("sonic failed to bind json, %v", err)
		}
		setGinHttpBody(c, input)

		handlers.RegisterUser(mock, createHash)(c)

		assertHTTPStatus(t, c, http.StatusCreated, w, "http status code should indicate that the user was created")
		field := mock.Users[0]

		assert.Equal(predefinedHash, field.PasswordHash)
		assert.Equal(username, field.Username)
		assert.Equal(name, *field.Name)
		assert.Equal(email, field.Email)
	})
	t.Run("user already exists", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		alreadyExistingUserName := "asdfasdf"
		alreadyExistinguser := models.User{
			ID: 1, Name: &alreadyExistingUserName, Username: "testuseralreadyexists", Email: "gregor@gregtech.eu",
		}
		mock.Users = append(mock.Users, alreadyExistinguser)

		input, err := sonic.Marshal(handlers.RegistrationInput{
			Username: alreadyExistinguser.Username, Email: alreadyExistinguser.Email, Password: "vverysecurepassword", Name: alreadyExistinguser.Name,
		})
		if err != nil {
			t.Fatalf("sonic failed to bind json: %v", err)
		}
		setGinHttpBody(c, input)

		handlers.RegisterUser(mock, createHash)(c)

		assertHTTPStatus(t, c, http.StatusConflict, w, "http status code should indicate that a user already exists")
		assert.Equal([]models.User{alreadyExistinguser}, mock.Users, "there should be just one user created")
	})
}
