package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	argon2id "github.com/gragorther/epigo/auth"
	"github.com/gragorther/epigo/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Env struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func initDB() *Env {
	env := &Env{}
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	env.DB = db
	env.DB.AutoMigrate(&models.User{})
	env.Logger = log.Default()
	return env
}

func (env *Env) registerUser(c *gin.Context) {
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

	var userFound models.User
	var emailFound models.User
	env.DB.Where("username=?", authInput.Username).Find(&userFound) // checks if user already exists
	env.DB.Where("email=?", authInput.Email).Find(&emailFound)
	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}
	if emailFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already used"})
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
	}

	env.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})

}

func main() {

	env := initDB()
	router := gin.Default()
	router.POST("/registerUser", env.registerUser)

	router.Run()
}
