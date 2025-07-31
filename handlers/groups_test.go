package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Err              error
	IsAuthorized     bool
	DeleteGroupCalls uint

	Groups     map[uint]models.Group
	Recipients []models.Recipient
}

func newMockDB(err error) *mockDB {
	m := mockDB{Err: err}
	// makes sure we don't get a panic because the map wasn't created
	m.Groups = make(map[uint]models.Group)

	return &m
}

func (m *mockDB) CreateGroupAndRecipientEmails(group *models.Group) error {

	// this gets the current length of the Groups map and sets the input group to the index at the uint of the length
	m.Groups[uint(len(m.Groups))] = *group
	return m.Err
}
func (m *mockDB) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	return m.IsAuthorized, m.Err
}
func (m *mockDB) DeleteGroupByID(id uint) error {
	m.DeleteGroupCalls += 1
	return m.Err
}

// sets up gin and returns a gin Context and an http response recorder
func setupGin() (*gin.Context, *httptest.ResponseRecorder) {

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

type invalidGroupInput struct {
	Recipients  []models.Recipient `json:"recipientsoopsietypo"` //simulates a typo in the json key
	Name        string             `json:"nameaaaeraser"`
	Description string             `json:"descriptionuu"`
}

func TestAddGroup(t *testing.T) {
	assert := assert.New(t)
	t.Run("with valid json input", func(t *testing.T) {
		c, w := setupGin()
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)

		mock := newMockDB(nil)
		jsonInput := handlers.GroupInput{
			Recipients: []models.Recipient{
				{Email: "gregor@gregtech.eu"},
				{Email: "test@gregtech.eu"},
			},
			Name:        "test name",
			Description: "test description",
		}
		jsonString, _ := sonic.Marshal(&jsonInput)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.AddGroup(mock)
		handler(c)
		assert.Equal(http.StatusOK, w.Code, "status code should indicate success")
		assert.Equal(jsonInput.Name, mock.Groups[0].Name)
		assert.Equal(jsonInput.Description, mock.Groups[0].Description)
		assert.Equal(fakeUser.ID, mock.Groups[0].UserID)
		assert.Equal(jsonInput.Recipients, mock.Groups[0].Recipients)
	})
	t.Run("with invalid json input", func(t *testing.T) {
		c, w := setupGin()
		mock := newMockDB(nil)
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)
		jsonInput, _ := sonic.Marshal(invalidGroupInput{
			Recipients: []models.Recipient{
				{Email: "uwu@gregtech.eu"},
			},
		})
		handler := handlers.AddGroup(mock)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonInput)),
		}

		handler(c)

		assert.Equal(http.StatusInternalServerError, w.Code, "status code should not indicate success")

	})
	t.Run("with invalid emails", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		mock := newMockDB(nil)
		handler := handlers.AddGroup(mock)
		jsonInput, _ := sonic.Marshal(handlers.GroupInput{
			Recipients: []models.Recipient{
				{Email: "asdf@"},
				{Email: "@email."},
			},
			Name:        "test group",
			Description: "test description",
		})
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonInput)),
		}
		handler(c)

		assert.Equal(http.StatusUnprocessableEntity, w.Code, "status code should indicate that the email is invalid")
	})
}

func TestDeleteGroup(t *testing.T) {
	assert := assert.New(t)
	t.Run("with valid input", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = true
		handler := handlers.DeleteGroup(mock)
		c.AddParam("id", "0")
		mock.Groups = map[uint]models.Group{
			0: {
				UserID: 1,
			},
		}

		handler(c)
		assert.Equal(http.StatusOK, w.Code, "status code should indicate success")

		assert.Equal(uint(1), mock.DeleteGroupCalls, "expected 1 delete group call")
	})
	t.Run("missing param", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = true
		handler := handlers.DeleteGroup(mock)

		handler(c)
		assert.Equal(http.StatusNotFound, w.Code, "group should not be found when there's no param")
		assert.Equal(uint(0), mock.DeleteGroupCalls, "expect 0 delete group calls when there's no param")
	})
}
