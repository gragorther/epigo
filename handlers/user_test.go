package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
)

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

		handlers.RegisterUser(mock)(c)

		assert.Equal(http.StatusCreated, w.Code, "http status code should indicate that the user was created")
		field := mock.Users[0]
		passMatch, err := argon2id.ComparePasswordAndHash(password, field.PasswordHash)
		if err != nil {
			t.Fatalf("failed to verify password hash match: %v", err)
		}
		if !passMatch {
			t.Errorf("Expected password hash of %v to match the one in the database", password)
		}

		assert.Equal(username, field.Username)
		assert.Equal(name, *field.Name)
		assert.Equal(email, field.Email)
	})
}
