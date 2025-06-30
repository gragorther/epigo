package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/db"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/util"
)

type UserHandler struct {
	U db.Users
}

type RegistrationInput struct {
	Username string  `json:"username" binding:"required"`
	Name     *string `json:"name"     binding:"required"`
	Email    string  `json:"email"    binding:"required,email"`
	Password string  `json:"password" binding:"required"`
}
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {

	var authInput RegistrationInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validEmail := util.ValidateEmail(authInput.Email)
	if !validEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrInvalidEmail.Error})
	}

	userExists, userExistsErr := h.U.CheckIfUserExistsByUsernameAndEmail(authInput.Username, authInput.Email)

	if userExistsErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check if the user already exists"})
		return
	}
	if userExists {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	passwordHash, err := argon2id.CreateHash(authInput.Password, argon2id.DefaultParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:     authInput.Username,
		PasswordHash: string(passwordHash),
		Email:        authInput.Email,
		Name:         authInput.Name,
	}

	if err := h.U.CreateUser(&user); err != nil {
		log.Printf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, authInput)

}
func (h *UserHandler) LoginUser(c *gin.Context) {
	var authInput LoginInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userExists, err := h.U.CheckIfUserExistsByUsername(authInput.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't check if user exists"})
		log.Print(err)
		return
	}
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": apperrors.ErrNotFound.Error()})
		return
	}
	userFound, err := h.U.GetUserByUsername(authInput.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "couldn't fetch user"})
		log.Printf("Couldn't fetch user: %v", err)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(authInput.Password, userFound.PasswordHash) //check hash
	if err != nil {
		log.Printf("Password hash checking error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error during password hash verification"})
		return
	}

	if !match {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate JWT token"})
		return
	}
	currentTime := time.Now()
	userFound.LastLogin = &currentTime
	if err := h.U.SaveUserData(userFound); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrServerError.Error()})
		log.Printf("failed to save user lastLogin: %v", err)
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})

}
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Retrieve the user object from the context
	userValue, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Type assertion to convert the interface{} to models.User
	user, ok := userValue.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrTypeConversionFailed.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"lastLogin": user.LastLogin,
		"name":      user.Name,
		"email":     user.Email,
	})
}

type setEmailIntervalInput struct {
	Cron string `json:"cron" binding:"required"`
}

func (h *UserHandler) SetEmailInterval(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(*models.User)
	var input setEmailIntervalInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": apperrors.ErrParsingFailed.Error()})
		return
	}
	err = h.U.UpdateUserInterval(user.ID, input.Cron)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("failed to set email interval: %v", err)
		return
	}
}
