package types

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
