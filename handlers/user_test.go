package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
func (m *mockDB) CheckIfUserExistsByUsername(username string) (bool, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockDB) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}
func (m *mockDB) EditUser(newUser *models.User) error {
	for _, user := range m.Users {
		if user.ID == newUser.ID {
			user = *newUser
			break
		}
	}
	return nil
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

func comparePasswordAndHash(password string, hash string) (bool, error) {

	hashedpass, _ := createHash(password, argon2id.DefaultParams)
	return hashedpass == predefinedHash, nil
}

func TestLoginUser(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()

		userPassword := "securepass123"
		hashedPass, _ := createHash(userPassword, argon2id.DefaultParams)
		// the user logging in
		userLoggingIn := models.User{
			ID: 1, Username: "test", PasswordHash: hashedPass,
		}
		mock.Users = append(mock.Users, userLoggingIn)

		input, err := sonic.Marshal(handlers.LoginInput{
			Username: userLoggingIn.Username, Password: userPassword,
		})
		if err != nil {
			t.Fatalf("sonic failed to marshal json: %v", err)
		}
		setGinHttpBody(c, input)

		jwtSecret := "secure jwt"
		handlers.LoginUser(mock, comparePasswordAndHash, jwtSecret)(c)

		assertHTTPStatus(t, c, http.StatusOK, w, "http status code should indicate success")

		var body struct {
			Token string `json:"token"`
		}

		if err := sonic.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Errorf("invalid output: %v, got error %v", body, err)
		}

		assert.NotEmpty(body.Token, "token should not be empty")

		token, err := jwt.Parse(body.Token, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil {
			t.Errorf("invalid token: %v", err)
		}

		if !token.Valid {
			t.Errorf("invalid token: %v", token.Raw)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("could not parse token claims")
		}

		assert.Equal(float64(userLoggingIn.ID), claims["id"])
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			t.Errorf("expired token (invalid time?) %v", claims["exp"])
			return
		}
	})
}
