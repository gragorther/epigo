package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *UserHandler) RegisterUser(c *gin.Context) {
	// var dto models.UserDTO
	// if err := c.BindJSON(&dto); err != nil {
	// 	return
	// }
	// passwordHash, err := argon2id.CreateHash(dto.Password, &argon2id.Params{Memory: 256 * 1024,
	// 	Iterations:  3,
	// 	Parallelism: 5,
	// 	SaltLength:  16,
	// 	KeyLength:   32})
	// if err != nil {
	// 	env.Logger.Printf("Password hashing error: %v", err)
	// }
	// newUser := models.User{
	// 	Username:     dto.Username,
	// 	Name:         dto.Name,
	// 	Email:        dto.Email,
	// 	PasswordHash: passwordHash,
	// 	LastLogin:    nil,
	// }
	// if err := env.DB.Create(&newUser).Error; err != nil {
	// 	env.Logger.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
	// 	return
	// }

	var authInput models.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists models.User
	res := h.DB.
		Where("username = ? OR email = ?", authInput.Username, authInput.Email).
		First(&exists)

	if res.Error == nil { // record found
		if exists.Username == authInput.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
		}
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
	}

	h.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})

}
