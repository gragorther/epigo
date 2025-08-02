package handlers_test

import (
	"bytes"
	"encoding/json"
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
	Err              error // the error mockDB methods will return
	IsAuthorized     bool
	DeleteGroupCalls uint
	UpdateGroupCalls uint
	FindGroupCalls   uint

	Groups       []models.Group
	LastMessages []models.LastMessage
	Recipients   []models.Recipient
}

func newMockDB(err error) *mockDB {
	m := mockDB{Err: err}

	return &m
}

func (m *mockDB) CreateGroupAndRecipientEmails(group *models.Group) error {

	// this gets the current length of the Groups map and sets the input group to the index at the uint of the length
	m.Groups = append(m.Groups, *group)
	return m.Err
}
func (m *mockDB) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	return m.IsAuthorized, m.Err
}
func (m *mockDB) DeleteGroupByID(id uint) error {
	m.DeleteGroupCalls += 1
	return m.Err
}

func (m *mockDB) FindGroupsAndRecipientsByUserID(userID uint) ([]models.Group, error) {
	m.FindGroupCalls += 1
	return m.Groups, m.Err
}
func (m *mockDB) UpdateGroup(group *models.Group) error {
	for i := range m.Groups {
		if m.Groups[i].ID == group.ID {
			m.Groups[i] = *group
			break
		}
	}

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
	Recipients  []models.APIRecipient `json:"recipientsoopsietypo"` //simulates a typo in the json key
	Name        string                `json:"nameaaaeraser"`
	Description string                `json:"descriptionuu"`
}

func TestAddGroup(t *testing.T) {
	assert := assert.New(t)
	t.Run("with valid json input", func(t *testing.T) {
		c, w := setupGin()
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)

		mock := newMockDB(nil)

		description := "test description"
		jsonInput := handlers.GroupInput{
			Recipients: []models.APIRecipient{
				{Email: "gregor@gregtech.eu"},
				{Email: "test@gregtech.eu"},
			},
			Name:        "test name",
			Description: &description,
		}
		jsonString, _ := sonic.Marshal(&jsonInput)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.AddGroup(mock)
		handler(c)
		assert.Equal(http.StatusOK, w.Code, "status code should indicate success")
		assert.Equal(jsonInput.Name, mock.Groups[0].Name)
		assert.Equal(*jsonInput.Description, *mock.Groups[0].Description)
		assert.Equal(fakeUser.ID, mock.Groups[0].UserID)
		for i, _ := range jsonInput.Recipients {
			assert.Equal(jsonInput.Recipients[i].Email, mock.Groups[0].Recipients[i].Email)
		}
	})
	t.Run("with invalid json input", func(t *testing.T) {
		c, w := setupGin()
		mock := newMockDB(nil)
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)
		jsonInput, _ := sonic.Marshal(invalidGroupInput{
			Recipients: []models.APIRecipient{
				{Email: "uwu@gregtech.eu"},
			},
		})
		handler := handlers.AddGroup(mock)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonInput)),
		}

		handler(c)

		assert.Equal(http.StatusUnprocessableEntity, w.Code, "status code should not indicate success")

	})
	t.Run("with invalid emails", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		mock := newMockDB(nil)
		handler := handlers.AddGroup(mock)

		description := "test description"
		jsonInput, _ := sonic.Marshal(handlers.GroupInput{
			Recipients: []models.APIRecipient{
				{Email: "asdf@"},
				{Email: "@email."},
			},
			Name:        "test group",
			Description: &description,
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
		mock.Groups = []models.Group{
			{
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
	t.Run("user does not own group", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = false

		handler := handlers.DeleteGroup(mock)
		c.AddParam("id", "0")

		handler(c)

		assert.Equal(http.StatusUnauthorized, w.Code, "user should be unauthorized to perform this action")
		assert.Equal(uint(0), mock.DeleteGroupCalls, "there should be no calls to delete the group because the user is unauthorized")
	})
}

func TestListGroups(t *testing.T) {
	assert := assert.New(t)

	c, w := setupGin()
	fakeUser := &models.User{ID: 1, Username: "test"}
	c.Set("currentUser", fakeUser)
	c.AddParam("id", "1")
	mock := newMockDB(nil)
	desc := "test desc"
	mock.Groups = []models.Group{
		{
			Name:        "test group 1",
			Description: &desc,
			Recipients: []models.Recipient{
				{APIRecipient: models.APIRecipient{Email: "gregor@gregtech.eu"}},
				{APIRecipient: models.APIRecipient{Email: "test@email.com"}},
			},
			LastMessages: []models.LastMessage{
				{Title: "test"},
			},
		},
	}
	handler := handlers.ListGroups(mock)

	handler(c)
	var unmarshaledBody []models.Group

	err := json.Unmarshal(w.Body.Bytes(), &unmarshaledBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal json: %v", err)
	}

	for i := range mock.Groups {
		assert.Equal(mock.Groups[i].Name, unmarshaledBody[i].Name)
		assert.Equal(desc, *unmarshaledBody[i].Description)
		assert.Equal(mock.Groups[i].LastMessages, unmarshaledBody[i].LastMessages)
		assert.Equal(mock.Groups[i].Recipients, unmarshaledBody[i].Recipients)
	}

	assert.Equal(http.StatusOK, w.Code, "http status code should indicate success")

	assert.Equal(mock.FindGroupCalls, uint(1), "expected 1 find group call")
}

func TestEditGroup(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid input", func(t *testing.T) {
		c, w := setupGin()
		fakeUser := &models.User{ID: 1, Username: "test"}
		c.Set("currentUser", fakeUser)
		c.AddParam("id", "0")
		mock := newMockDB(nil)
		mock.IsAuthorized = true
		// the constants to be used in the group.
		const groupID uint = 0
		const unchangedGroupName string = "not yet changed"
		var unchangedGroupDesc string = "unchanged group desc"

		// this can't be const for some reason
		var unchangedGroupRecipients []models.Recipient = []models.Recipient{
			{GroupID: groupID, APIRecipient: models.APIRecipient{Email: "gregor@gregtech.eu"}},
			{GroupID: groupID, APIRecipient: models.APIRecipient{Email: "gregor@gregtech.eu"}},
		}
		mock.Groups = []models.Group{

			{Name: unchangedGroupName, Description: &unchangedGroupDesc, Recipients: unchangedGroupRecipients, ID: groupID},
		}

		// the new names of fields in the group which will be then asserted against to see if handler actually changed anything
		const newGroupName string = "new name :3"
		var newGroupDesc string = "new desc"
		newRecipients := []models.APIRecipient{
			{Email: "test@thing.com"},
		}
		jsonString, _ := sonic.Marshal(&handlers.GroupInput{

			Name:        newGroupName,
			Description: &newGroupDesc,
			Recipients:  newRecipients,
		})
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.EditGroup(mock)
		handler(c)

		assert.Equal(http.StatusOK, w.Code, "http status code should indicate success")

		newRecipient := newRecipients[0]
		recipient := mock.Groups[0].Recipients[0]
		assert.Equal(newGroupName, mock.Groups[0].Name, "the group name in the in-memory database should match the one sent to the API")
		assert.Equal(newGroupDesc, *mock.Groups[0].Description, "the group descripton in the in memory db should match the one sent to the api")
		assert.Equal(newRecipient.Email, recipient.Email, "the recipients sent to the api should match the ones in the in memory db")

	})
}
