package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/email"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/models"
)

type UserUserStore interface {
	CreateUser(*models.User) error
	CheckIfUserExistsByUsername(username string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	SaveUserData(*models.User) error
	UpdateUserInterval(userID uint, cron string) error
	CheckIfUserExistsByUsernameAndEmail(username string, email string) (bool, error)
}

type UserHandler struct {
	u UserUserStore
}

func NewUserHandler(u UserUserStore) *UserHandler {
	return &UserHandler{u: u}
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
		c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
		return
	}
	validEmail := email.Validate(authInput.Email)
	if !validEmail {
		c.AbortWithError(http.StatusBadRequest, apperrors.ErrInvalidEmail)
		return
	}

	userExists, userExistsErr := h.u.CheckIfUserExistsByUsernameAndEmail(authInput.Username, authInput.Email)

	if userExistsErr != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
		return
	}
	if userExists {
		c.AbortWithError(http.StatusConflict, apperrors.ErrUserAlreadyExists)
		return
	}

	passwordHash, err := argon2id.CreateHash(authInput.Password, argon2id.DefaultParams)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrHashingFailed)
		return
	}

	user := models.User{
		Username:     authInput.Username,
		PasswordHash: string(passwordHash),
		Email:        authInput.Email,
		Name:         authInput.Name,
	}

	if err := h.u.CreateUser(&user); err != nil {
		log.Printf("failed to create user: %v", err)
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
		return
	}

	c.JSON(http.StatusOK, authInput)

}
func (h *UserHandler) LoginUser(c *gin.Context) {
	var authInput LoginInput
	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
		return
	}

	userExists, err := h.u.CheckIfUserExistsByUsername(authInput.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrUserNotFound)
		log.Print(err)
		return
	}
	if !userExists {
		c.AbortWithError(http.StatusNotFound, apperrors.ErrUserNotFound)
		return
	}
	userFound, err := h.u.GetUserByUsername(authInput.Username)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, apperrors.ErrUserNotFound)
		log.Printf("Couldn't fetch user: %v", err)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(authInput.Password, userFound.PasswordHash) //check hash
	if err != nil {
		log.Printf("Password hash checking error: %v", err)
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrHashCheckFailed)
		return
	}

	if !match {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidPassword)
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrJWTCreationError)
		return
	}
	currentTime := time.Now()
	userFound.LastLogin = &currentTime
	if err := h.u.SaveUserData(userFound); err != nil {
		c.Error(apperrors.ErrCreationOfObjectFailed)
		log.Printf("failed to save user lastLogin: %v", err)
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})

}
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// Retrieve the user object from the context
	userValue, _ := c.Get("currentUser")

	// Type assertion to convert the interface{} to models.User
	user, ok := userValue.(*models.User)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
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
		c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
		return
	}
	err = h.u.UpdateUserInterval(user.ID, input.Cron)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
		log.Printf("failed to set email interval: %v", err)
		return
	}
}
