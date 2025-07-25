package stubdb

type AuthDB struct{}

func (a AuthDB) CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error) {
	return true, nil
}
func (a AuthDB) CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error) {
	return true, nil
}
