package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	env.DB.Create(&newUser)
	c.JSON(http.StatusCreated, newUser)
}

func main() {
	env := initDB()
	router := gin.Default()
	router.POST("/registerUser", env.registerUser)

	router.Run()
}
