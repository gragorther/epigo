package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	argon2id "github.com/gragorther/epigo/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Env struct {
	DB     *gorm.DB
	Logger *log.Logger
}

type User struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	//	Country      string //should probably be a foreign key of another table
	LastLogin *time.Time `json:"lastLogin"`
}
type UserDTO struct { // prevents client from modifying everything in the users table
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func initDB() *Env {
	env := &Env{}
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	env.DB = db
	env.DB.AutoMigrate(&User{})
	env.Logger = log.Default()
	return env
}

func (env *Env) registerUser(c *gin.Context) {
	var dto UserDTO
	if err := c.BindJSON(&dto); err != nil {
		return
	}
	passwordHash, err := argon2id.CreateHash(dto.Password, &argon2id.Params{Memory: 256 * 1024,
		Iterations:  3,
		Parallelism: 5,
		SaltLength:  16,
		KeyLength:   32})
	if err != nil {
		env.Logger.Printf("Password hashing error: %v", err)
	}
	newUser := User{
		Username:     dto.Username,
		Name:         dto.Name,
		Email:        dto.Email,
		PasswordHash: passwordHash,
		LastLogin:    nil,
	}
	if err := env.DB.Create(&newUser).Error; err != nil {
		env.Logger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}

}

func main() {

	env := initDB()
	router := gin.Default()
	router.POST("/registerUser", env.registerUser)

	router.Run()
}
