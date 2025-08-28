package mock

import (
	"context"

	"github.com/gragorther/epigo/models"
)

func (m *mockDB) CreateLastMessage(ctx context.Context, lastMessage *models.LastMessage) error {
	m.LastMessages = append(m.LastMessages, *lastMessage)
	return m.Err
}
func (m *mockDB) FindLastMessagesByUserID(userID uint) ([]models.LastMessage, error) {
	var output []models.LastMessage

	for _, message := range m.LastMessages {
		if message.UserID == userID {
			output = append(output, message)
		}
	}

	return output, m.Err
}
func (m *mockDB) UpdateLastMessage(ctx context.Context, newMessage models.LastMessage) error {
	for i, message := range m.LastMessages {
		if message.ID == newMessage.ID {
			m.LastMessages[i] = newMessage
		}
	}
	return m.Err
}

func (m *mockDB) CheckUserAuthorizationForLastMessage(ctx context.Context, messageID uint, userID uint) (bool, error) {
	for _, message := range m.LastMessages {
		if message.ID == messageID {
			if message.UserID == userID {
				return true, m.Err
			}
		}
	}
	return false, m.Err
}
