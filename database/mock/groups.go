package mock

import (
	"context"

	"github.com/gragorther/epigo/models"
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
	Users        []models.User
	Profiles     []models.Profile
}

func NewMockDB() *mockDB {
	m := mockDB{}

	return &m
}

func (m *mockDB) CreateGroup(group *models.Group) error {

	// this gets the current length of the Groups map and sets the input group to the index at the uint of the length
	m.Groups = append(m.Groups, *group)
	return m.Err
}
func (m *mockDB) CheckUserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error) {

	return m.IsAuthorized, m.Err
}
func (m *mockDB) DeleteGroupByID(ctx context.Context, id uint) error {
	m.DeleteGroupCalls += 1
	return m.Err
}

func (m *mockDB) FindGroupsAndRecipientsByUserID(ctx context.Context, userID uint) ([]models.Group, error) {
	m.FindGroupCalls += 1
	return m.Groups, m.Err
}
func (m *mockDB) UpdateGroup(ctx context.Context, group models.Group) error {
	for i := range m.Groups {
		if m.Groups[i].ID == group.ID {
			m.Groups[i] = group
			break
		}
	}

	return m.Err
}
