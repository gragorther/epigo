package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/cron"
	"github.com/gragorther/epigo/database/db"
	dbHandlers "github.com/gragorther/epigo/database/db"
	ginctx "github.com/gragorther/epigo/handlers/context"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/tokens"
	"github.com/guregu/null/v6"
)

type RegistrationInput struct {
	Username string      `json:"username" binding:"required"`
	Name     null.String `json:"name"`
	Password string      `json:"password" binding:"required"`
	Token    string      `json:"token" binding:"required"`
}
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// emailVerificationRoute is the route with a token query parameter that is used for verifying the user email
//
// this also takes a `token` query parameter, which is the email JWT token
func Register(db interface {
	CheckIfUserExistsByUsernameAndEmail(ctx context.Context, username string, email string) (bool, error)
}, queue interface {
	CreateUser(db.CreateUserInput) error
}, createHash func(string, *argon2id.Params) (string, error), parseEmailVerificationToken tokens.ParseEmailVerificationFunc,
) gin.HandlerFunc {
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

		userExists, err := db.CheckIfUserExistsByUsernameAndEmail(c, authInput.Username, userEmail)
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

		if err := queue.CreateUser(dbHandlers.CreateUserInput{Username: authInput.Username, Email: userEmail, Name: authInput.Name, PasswordHash: passwordHash}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
			return
		}
		c.Status(http.StatusCreated)
	}
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(db interface {
	UserIDAndPasswordHashByUsername(ctx context.Context, username string) (user db.UserIDAndPasswordHash, err error)
	CheckIfUserExistsByUsername(ctx context.Context, username string) (bool, error)
}, comparePasswordAndHash func(password string, hash string) (match bool, err error), createUserAuthToken tokens.CreateUserAuthFunc,
) gin.HandlerFunc {
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
		userFound, err := db.UserIDAndPasswordHashByUsername(c, authInput.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user by username: %w", err))
			return
		}
		match, err := comparePasswordAndHash(authInput.Password, userFound.PasswordHash) // check hash
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

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
		})
	}
}

type UpdateMaxSentEmailsInput struct {
	MaxSentEmails uint `json:"maxSentEmails" binding:"required"`
}

func UpdateMaxSentEmails(queue interface {
	SetUserMaxSentEmails(userID uint, maxSentEmails uint) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input UpdateMaxSentEmailsInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := queue.SetUserMaxSentEmails(userID, input.MaxSentEmails); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

type GetUserDataOutput struct {
	Username string `json:"username,omitzero"`
	Name     string `json:"name,omitzero"`
	Email    string `json:"email,omitzero"`
}

func GetData(db interface {
	UserByID(context.Context, uint) (db.User, error)
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user object from the context
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user profile: %w", err))
			return
		}
		user, err := db.UserByID(c, userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user from ID: %w", err))
			return
		}

		c.JSON(http.StatusOK, GetUserDataOutput{
			Username: user.Username,
			Name:     user.Name.String,
			Email:    user.Email,
		})
	}
}

type setEmailIntervalInput struct {
	Cron string `json:"cron" binding:"required"`
}

func SetEmailInterval(queue interface {
	UpdateUserInterval(uint, string) error
}, minDurationBetweenEmails time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input setEmailIntervalInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to bind json while setting user email interval: %w", err))
			return
		}
		minDurationBetweenTicks, err := cron.MinDurationBetweenCronTicks(input.Cron, 0)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if minDurationBetweenTicks < minDurationBetweenEmails {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if err := queue.UpdateUserInterval(userID, input.Cron); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update user interval: %w", err))
			return
		}
	}
}

type EmailVerificationInput struct {
	Email string `json:"email" binding:"required,email"`
}

// here, the user enters their email, gets sent a registration link to their email, and continues registration from there
func VerifyEmail(queue interface {
	SendVerificationEmail(email string) error
}, db interface {
	CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error)
},
) gin.HandlerFunc {
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
		if err := queue.SendVerificationEmail(input.Email); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
