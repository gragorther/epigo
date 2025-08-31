package mock

import (
	"context"

	"github.com/gragorther/epigo/models"
)

func (m *mockDB) CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error) {

	for _, user := range m.Users {
		if user.Username == username && user.Email == email {
			return true, nil
		}
	}
	return false, nil
}
func (m *mockDB) CheckIfUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockDB) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return models.User{}, nil
}

func (m *mockDB) GetUserByID(ctx context.Context, userID uint) (models.User, error) {
	for _, user := range m.Users {
		if user.ID == userID {
			return user, nil
		}
	}
	return models.User{}, nil
}

func (m *mockDB) EditUser(ctx context.Context, newUser models.User) error {
	for i, user := range m.Users {
		if user.ID == newUser.ID {
			m.Users[i] = newUser
			break
		}
	}
	return nil
}

func (m *mockDB) CreateUser(ctx context.Context, newUser *models.User) error {
	m.Users = append(m.Users, *newUser)
	return nil
}

func (m *mockDB) UpdateProfile(_ context.Context, newProfile models.Profile) error {
	for i, profile := range m.Profiles {
		if profile.UserID == newProfile.UserID {
			m.Profiles[i] = newProfile
			break
		}
	}
	return nil
}

func (m *mockDB) CreateProfile(ctx context.Context, newProfile *models.Profile) error {
	m.Profiles = append(m.Profiles, *newProfile)
	return nil
}

func (m *mockDB) CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}
