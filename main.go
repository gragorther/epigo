package main

import (
	"fmt"
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
	ID           uint
	Username     string
	Name         string
	Email        string
	PasswordHash string
	//	Country      string //should probably be a foreign key of another table
	LastLogin *time.Time
}

func initDB() {
	env := &Env{}
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error opening database: %v", err)
		return
	}
	env.DB = db // sets the DB env to the database connector set up above
	env.DB.AutoMigrate(&User{})

}

func registerUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	c.JSON(http.StatusCreated, newUser)

}

func main() {
	initDB()
	router := gin.Default()
	router.POST("/registerUser", registerUser)
}
