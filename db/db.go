package db

type DB interface {
	Auth
	Groups
	Messages
	Users
}

type DBHandler struct {
	Auth     Auth
	Messages Messages
	Users    Users
	Groups   Groups
}
