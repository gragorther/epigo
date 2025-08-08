package handlers_test

import (
	"context"
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
	for i, user := range m.Users {
		if user.ID == newUser.ID {
			m.Users[i] = *newUser
			break
		}
	}
	return nil
}
func (m *mockDB) CreateUser(newUser *models.User) error {
	m.Users = append(m.Users, *newUser)
	return nil
}
func (m *mockDB) UpdateProfile(_ context.Context, newProfile models.Profile) error {
	for i, profile := range m.Profiles {
		if profile.UserID == newProfile.UserID {
			m.Profiles[i] = newProfile
			break
		}
	}
	return nil
}
func (m *mockDB) CreateProfile(newProfile *models.Profile) error {
	m.Profiles = append(m.Profiles, *newProfile)
	return nil
}

// a predefined hash so tests are much faster to run because they don't have to *actually* hash the password
const hashSuffix string = "this is hashed"

// a mock of argon2id.CreateHash
func createHash(password string, params *argon2id.Params) (string, error) {
	return password + hashSuffix, nil
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

		hash, _ := createHash(password, argon2id.DefaultParams)
		assert.Equal(hash, field.PasswordHash)
		assert.Equal(username, field.Username)
		assert.Equal(name, *field.Profile.Name)
		assert.Equal(email, field.Email)
	})
	t.Run("user already exists", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		alreadyExistingUserName := "asdfasdf"
		alreadyExistinguser := models.User{
			ID: 1, Profile: &models.Profile{Name: &alreadyExistingUserName}, Username: "testuseralreadyexists", Email: "gregor@gregtech.eu",
		}
		mock.Users = append(mock.Users, alreadyExistinguser)

		input, err := sonic.Marshal(handlers.RegistrationInput{
			Username: alreadyExistinguser.Username, Email: alreadyExistinguser.Email, Password: "vverysecurepassword", Name: alreadyExistinguser.Profile.Name,
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
	hashedPass, _ := createHash(password, argon2id.DefaultParams)
	return hashedPass == hash, nil
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
			t.Fatalf("invalid token: %v", err)
		}

		if !token.Valid {
			t.Fatalf("invalid token: %v", token.Raw)
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
	t.Run("invalid password", func(t *testing.T) {
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
			Username: userLoggingIn.Username, Password: "invalid password oopsie",
		})
		if err != nil {
			t.Fatalf("sonic failed to marshal json: %v", err)
		}
		setGinHttpBody(c, input)
		jwtSecret := "secure jwt !!"

		handlers.LoginUser(mock, comparePasswordAndHash, jwtSecret)(c)

		assertHTTPStatus(t, c, http.StatusUnauthorized, w, "user should be unauthorized")

		assert.Empty(w.Body.Bytes())

	})
	t.Run("user not found", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)

		mock := newMockDB()
		input, err := sonic.Marshal(handlers.LoginInput{
			Username: "iDontExist", Password: "invalid password oopsie",
		})
		if err != nil {
			t.Fatalf("sonic failed to marshal json: %v", err)
		}
		setGinHttpBody(c, input)
		jwtSecret := "secure jwt !!"

		handlers.LoginUser(mock, comparePasswordAndHash, jwtSecret)(c)
		assertHTTPStatus(t, c, http.StatusNotFound, w, "http status code should indicate that the user was not found")
		assert.Empty(w.Body.Bytes(), "the response body should be empty because the user was not found")
	})
}

func TestGetUserData(t *testing.T) {
	c, w, assert := setupHandlerTest(t)
	userName := "myname"
	lastLoginTime := time.Date(2025, time.January, 1, 12, 12, 12, 12, time.UTC)
	user := &models.User{
		Profile:   &models.Profile{Name: &userName},
		Username:  "myusername",
		LastLogin: &lastLoginTime,
		Email:     "eeemail@gregtech.eu",
	}
	handlers.SetUser(c, user)
	handlers.GetUserData()(c)

	var output struct {
		Username  string    `json:"username"`
		LastLogin time.Time `json:"lastLogin"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
	}
	assertHTTPStatus(t, c, http.StatusOK, w, "http status should indicate success")

	sonic.Unmarshal(w.Body.Bytes(), &output)
	assert.Equal(userName, output.Name)
	assert.Equal(user.Username, output.Username)
	assert.Equal(lastLoginTime, output.LastLogin)
	assert.Equal(user.Email, output.Email)
}
func TestCreateProfile(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		handlers.SetUser(c, &models.User{
			ID:       1,
			Username: "username",
		})
		profileInput := handlers.ProfileInput{
			Name: "newname",
		}
		input, err := sonic.Marshal(profileInput)
		if err != nil {
			t.Fatalf("sonic failed to bind json: %v", err)
		}
		setGinHttpBody(c, input)
		handlers.CreateProfile(mock)(c)

		assertHTTPStatus(t, c, http.StatusCreated, w, "http status should indicate that the profile was created")
		assert.Equal(*mock.Profiles[0].Name, profileInput.Name)
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		c, w, assert := setupHandlerTest(t)
		mock := newMockDB()
		profileInput := handlers.ProfileInput{
			Name: "uwu",
		}
		input, err := sonic.Marshal(profileInput)
		if err != nil {
			t.Fatalf("sonic failed to marshal json: %v", err)
		}

		oldName := "oldn naem"
		mock.CreateProfile(&models.Profile{
			Name:   &oldName,
			UserID: 1,
		})
		currentUser := models.User{
			Username: "bob", ID: 1,
		}
		handlers.SetUser(c, &currentUser)

		setGinHttpBody(c, input)
		handlers.UpdateProfile(mock)(c)

		assertHTTPStatus(t, c, http.StatusNoContent, w, "http status code should indicate that the profile was updated")
		assert.Equal(profileInput.Name, *mock.Profiles[0].Name)
	})
}
