package models

import (
	"time"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type Group struct {
	ID           uint           `json:"ID,omitzero" gorm:"primarykey"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	UserID       uint           `json:"-"`
	Name         string         `json:"name"`
	Description  *string        `json:"description,omitempty,omitzero"`
	Recipients   *[]Recipient   `json:"recipients"`
	LastMessages *[]LastMessage `json:"lastMessages" gorm:"many2many:group_last_messages;"`
}

// this, in combination with the omitzero tag, allows us to disable
func (g *Group) UnmarshalJSON(data []byte) error {
	type tempGroup struct {
		Name         string        `json:"name"`
		Description  string        `json:"description"`
		Recipients   []Recipient   `json:"recipients"`
		LastMessages []LastMessage `json:"lastMessages" gorm:"many2many:group_last_messages;"`
	}

	var temp tempGroup
	err := sonic.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	g.Name = temp.Name
	g.Description = &temp.Description
	g.Recipients = &temp.Recipients
	g.LastMessages = &temp.LastMessages
	return nil
}

type Recipient struct {
	gorm.Model
	GroupID uint   `json:"groupID"` // group of the email
	Email   string `json:"email"`
}
