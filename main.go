package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db/initializers"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
)

func main() {
	db, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	r := gin.Default()
	userHandler := handlers.NewUserHandler(db)
	r.POST("/register", userHandler.RegisterUser)

	r.Run()
}
