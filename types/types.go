package types

import "time"

type LastMessageOut struct {
	ID      uint   `gorm:"primarykey"`
	Title   string `json:"title"`
	Groups  []uint `json:"groups" gorm:"many2many:group_last_messages;"`
	Content string `json:"content"`
}
type UserIntervalsOutput struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}
type GroupWithEmails struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	RecipientEmails []string  `json:"recipientEmails"`
}
