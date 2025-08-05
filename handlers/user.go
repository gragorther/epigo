package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gragorther/epigo/email"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
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
}) gin.HandlerFunc {
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

		passwordHash, err := argon2id.CreateHash(authInput.Password, argon2id.DefaultParams)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to register user: %w", err))
			return
		}

		user := models.User{
			Username:     authInput.Username,
			PasswordHash: string(passwordHash),
			Email:        authInput.Email,
			Name:         authInput.Name,
		}

		if err := db.CreateUser(&user); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create user: %w", err))
			return
		}

		c.Status(http.StatusCreated)
	}
}

func LoginUser(db interface {
	CheckIfUserExistsByUsername(username string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	SaveUserData(*models.User) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authInput LoginInput
		if err := c.ShouldBindJSON(&authInput); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("Failed to parse login input: %w", err))
			return
		}

		userExists, err := db.CheckIfUserExistsByUsername(authInput.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to check if users exists by username: %w", err))
			return
		}
		if !userExists {
			c.AbortWithError(http.StatusNotFound, errors.New("user not found"))
			return
		}
		userFound, err := db.GetUserByUsername(authInput.Username)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to get user by username: %w", err))
			return
		}
		match, err := argon2id.ComparePasswordAndHash(authInput.Password, userFound.PasswordHash) //check hash
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to compare password and hash: %w", err))
			return
		}

		if !match {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  userFound.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

		token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to generate JWT token: %w", err))
			return
		}
		currentTime := time.Now()
		userFound.LastLogin = &currentTime
		if err := db.SaveUserData(userFound); err != nil {
			c.Error(fmt.Errorf("Failed to store user last login: %w", err))
		}
		c.JSON(200, gin.H{
			"token": token,
		})
	}
}
func GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user object from the context
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to get user profile: %w", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  user.Username,
			"lastLogin": user.LastLogin,
			"name":      user.Name,
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
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("Failed to bind json while setting user email interval: %w", err))
			return
		}
		err = db.UpdateUserInterval(user.ID, input.Cron)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to update user interval: %w", err))
			return
		}
	}
}
