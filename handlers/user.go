package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	argon2id "github.com/gragorther/epigo/auth"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

type RegistrationInput struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

	var foundUsers int64
	res := h.DB.Model(&models.User{}).
		Where("username = ? OR email = ?", authInput.Username, authInput.Email).Count(&foundUsers)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check if user exists"})
		log.Printf("Couldn't check if user exists: %v", res.Error)
		return
	}

	if foundUsers > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
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

	h.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": authInput})

}
func (h *UserHandler) LoginUser(c *gin.Context) {
	var authInput LoginInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	h.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
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
	h.DB.Save(&userFound)
	c.JSON(200, gin.H{
		"token": token,
	})

}
func GetUserProfile(c *gin.Context) {
	// Retrieve the user object from the context
	userValue, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Type assertion to convert the interface{} to models.User
	user, ok := userValue.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":  user.Username,
		"lastLogin": user.LastLogin,
		"name":      user.Name,
		"email":     user.Email,
	})
}
