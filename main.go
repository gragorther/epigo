package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db/initializers"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/models"
)

func main() {
	db, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.RecipientEmail{}, &models.Group{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	r := gin.Default()
	userHandler := handlers.NewUserHandler(db)
	authHandler := middlewares.NewAuthHandler(db)
	groupHandler := handlers.NewGroupHandler(db)

	// user stuff
	r.POST("/user/register", userHandler.RegisterUser)
	r.POST("/user/login", userHandler.LoginUser)
	r.GET("/user/profile", authHandler.CheckAuth, handlers.GetUserProfile)

	// groups
	r.DELETE("/user/groups/delete/:id", authHandler.CheckAuth, groupHandler.DeleteGroup)
	r.POST("/user/groups/add", authHandler.CheckAuth, groupHandler.AddGroup)
	r.GET("/user/groups", authHandler.CheckAuth, groupHandler.ListGroups) // list groups
	r.PATCH("/user/groups/edit/:id", authHandler.CheckAuth, groupHandler.EditGroup)

	r.Run()
}
