package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adhocore/gronx"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/asynq/tasks"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/tokens"
)

type RegistrationInput struct {
	Username string  `json:"username" binding:"required"`
	Name     *string `json:"name"`
	Password string  `json:"password" binding:"required"`
	Token    string  `json:"token" binding:"required"`
}
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// emailVerificationRoute is the route with a token query parameter that is used for verifying the user email
//
// this also takes a `token` query parameter, which is the email JWT token
func RegisterUser(db interface {
	CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error)
	CreateUser(context.Context, *models.User) error
}, createHash func(string, *argon2id.Params) (string, error), parseEmailVerificationToken tokens.ParseEmailVerificationFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		var authInput RegistrationInput

		if err := c.ShouldBindJSON(&authInput); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind register user JSON: %w", err))
			return
		}
		userEmail, err := parseEmailVerificationToken(authInput.Token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to parse user email verification token, which was acquired from the `token` query param: %w", err))
			return
		}

		userExists, err := db.CheckIfUserExistsByUsernameAndEmail(authInput.Username, userEmail)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check if user exists by username and email: %w", err))
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
			PasswordHash: passwordHash,
			Email:        userEmail,
			Profile:      &models.Profile{Name: authInput.Name},
		}

		if err := db.CreateUser(c, &user); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
			return
		}
		c.Status(http.StatusCreated)
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginUser(db interface {
	CheckIfUserExistsByUsername(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	EditUser(context.Context, models.User) error
}, comparePasswordAndHash func(password string, hash string) (match bool, err error), createUserAuthToken tokens.CreateUserAuthFunc) gin.HandlerFunc {
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

		token, err := createUserAuthToken(userFound.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to generate JWT token: %w", err))
			return
		}
		currentTime := time.Now()
		userFound.LastLogin = &currentTime
		if err := db.EditUser(c, userFound); err != nil {
			c.Error(fmt.Errorf("failed to store user last login: %w", err))
		}
		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}

type GetUserDataOutput struct {
	Username  string    `json:"username,omitzero"`
	LastLogin time.Time `json:"lastLogin,omitzero"`
	Name      string    `json:"name,omitzero"`
	Email     string    `json:"email,omitzero"`
}

func GetUserData(db interface {
	GetUserByID(ctx context.Context, ID uint) (models.User, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user object from the context
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user profile: %w", err))
			return
		}
		user, err := db.GetUserByID(c, userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user from ID: %w", err))
			return
		}
		var lastLogin time.Time
		if user.LastLogin != nil {
			lastLogin = *user.LastLogin
		}
		var name string
		if user.Profile != nil && user.Profile.Name != nil {
			name = *user.Profile.Name
		}

		c.JSON(http.StatusOK, GetUserDataOutput{
			Username:  user.Username,
			LastLogin: lastLogin,
			Name:      name,
			Email:     user.Email,
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
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input setEmailIntervalInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to bind json while setting user email interval: %w", err))
			return
		}
		if !gronx.IsValid(input.Cron) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := db.UpdateUserInterval(userID, input.Cron); err != nil {
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
		userID, err := GetUserIDFromContext(c)
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
			UserID: userID,
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
		userID, err := GetUserIDFromContext(c)
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
			UserID: userID,
			Name:   &input.Name,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update profile: %w", err))
		}
		c.Status(http.StatusNoContent)
	}
}

type EmailVerificationInput struct {
	Email string `json:"email" binding:"required,email"`
}

// here, the user enters their email, gets sent a registration link to their email, and continues registration from there
func VerifyEmail(enqueueTask tasks.TaskEnqueuer, db interface {
	CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input EmailVerificationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
		exists, err := db.CheckIfUserExistsByEmail(c, input.Email)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if exists {
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		task, err := tasks.NewVerificationEmailTask(input.Email)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if _, err := enqueueTask(c, task); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
