package db

type Auth interface {
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error)
}
