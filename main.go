package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db/initializers"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/models"
)

func main() {
	fmt.Printf("JWT secret is: %v\n", os.Getenv("JWT_SECRET"))
	db, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.SendToGroup{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	r := gin.Default()
	userHandler := handlers.NewUserHandler(db)
	authHandler := middlewares.NewAuthHandler(db)
	r.POST("/register", userHandler.RegisterUser)
	r.POST("/login", userHandler.LoginUser)
	r.GET("/user/profile", authHandler.CheckAuth, handlers.GetUserProfile)
	r.Run()
}
