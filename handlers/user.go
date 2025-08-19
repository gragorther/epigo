package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/email"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/tokens"
)

type RegistrationInput struct {
	Username string  `json:"username" binding:"required"`
	Name     *string `json:"name"`
	Email    string  `json:"email"    binding:"required,email"`
	Password string  `json:"password" binding:"required"`
}
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(db interface {
	CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error)
	CreateUser(*models.User) error
}, createHash func(string, *argon2id.Params) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		var authInput RegistrationInput

		if err := c.ShouldBindJSON(&authInput); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind register user JSON: %w", err))
			return
		}
		validEmail := email.Validate(authInput.Email)
		if !validEmail {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		userExists, err := db.CheckIfUserExistsByUsernameAndEmail(authInput.Username, authInput.Email)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err))
			return
		}
		if userExists {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		passwordHash, err := createHash(authInput.Password, argon2id.DefaultParams)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err))
			return
		}

		user := models.User{
			Username:     authInput.Username,
			PasswordHash: string(passwordHash),
			Email:        authInput.Email,
			Profile:      &models.Profile{Name: authInput.Name},
		}

		if err := db.CreateUser(&user); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
			return
		}
		c.Status(http.StatusCreated)
	}
}

func LoginUser(db interface {
	CheckIfUserExistsByUsername(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	EditUser(context.Context, models.User) error
}, comparePasswordAndHash func(password string, hash string) (match bool, err error), jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authInput LoginInput
		if err := c.ShouldBindJSON(&authInput); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to parse login input: %w", err))
			return
		}

		userExists, err := db.CheckIfUserExistsByUsername(c, authInput.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check if users exists by username: %w", err))
			return
		}
		if !userExists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		userFound, err := db.GetUserByUsername(c, authInput.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user by username: %w", err))
			return
		}
		match, err := comparePasswordAndHash(authInput.Password, userFound.PasswordHash) //check hash
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to compare password and hash: %w", err))
			return
		}

		if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := tokens.CreateUserAuth(jwtSecret, userFound.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to generate JWT token: %w", err))
			return
		}
		currentTime := time.Now()
		userFound.LastLogin = &currentTime
		if err := db.EditUser(c, userFound); err != nil {
			c.Error(fmt.Errorf("failed to store user last login: %w", err))
		}
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
func GetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user object from the context
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user profile: %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  user.Username,
			"lastLogin": user.LastLogin,
			"name":      user.Profile.Name,
			"email":     user.Email,
		})
	}
}

type setEmailIntervalInput struct {
	Cron string `json:"cron" binding:"required"`
}

func SetEmailInterval(db interface {
	UpdateUserInterval(userID uint, cron string) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		var input setEmailIntervalInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind json while setting user email interval: %w", err))
			return
		}
		err = db.UpdateUserInterval(user.ID, input.Cron)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update user interval: %w", err))
			return
		}
	}
}

type ProfileInput struct {
	Name string `json:"name,omitempty"`
}

func CreateProfile(db interface {
	CreateProfile(context.Context, *models.Profile) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to get user from context: %w", err))
			return
		}
		var input ProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := db.CreateProfile(c, &models.Profile{
			UserID: user.ID,
			Name:   &input.Name,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.Status(http.StatusCreated)
	}
}

func UpdateProfile(db interface {
	UpdateProfile(context.Context, models.Profile) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user from context: %v", err))
			return
		}
		var input ProfileInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		//check for error
		if err := db.UpdateProfile(c, models.Profile{
			UserID: user.ID,
			Name:   &input.Name,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update profile: %w", err))
		}
		c.Status(http.StatusNoContent)
	}
}
