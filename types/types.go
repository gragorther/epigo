package types

type UserIntervalsOutput struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `json:"email" gorm:"unique"`
	EmailCron string `json:"emailCron"`
}
