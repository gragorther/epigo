package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/initializers"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/models"
)

func main() {
	dbconn, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = dbconn.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.RecipientEmail{}, &models.Group{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	dbHandler := db.NewDBHandler(dbconn)

	r := gin.Default()
	userHandler := handlers.NewUserHandler(dbconn)
	authHandler := middlewares.NewAuthHandler(dbconn)
	groupHandler := handlers.NewGroupHandler(dbHandler)
	messageHandler := handlers.NewMessageHandler(dbHandler)

	// user stuff
	r.POST("/user/register", userHandler.RegisterUser)
	r.POST("/user/login", userHandler.LoginUser)
	r.GET("/user/profile", authHandler.CheckAuth, handlers.GetUserProfile)

	// groups
	r.DELETE("/user/groups/delete/:id", authHandler.CheckAuth, groupHandler.DeleteGroup)
	r.POST("/user/groups/add", authHandler.CheckAuth, groupHandler.AddGroup)
	r.GET("/user/groups", authHandler.CheckAuth, groupHandler.ListGroups) // list groups
	r.PATCH("/user/groups/edit/:id", authHandler.CheckAuth, groupHandler.EditGroup)

	// lastMessages
	r.POST("/user/lastMessages/add", authHandler.CheckAuth, messageHandler.AddLastMessage)
	r.GET("/user/lastMessages", authHandler.CheckAuth, messageHandler.ListLastMessages)
	r.PATCH("/user/lastMessages/edit/:id", authHandler.CheckAuth, messageHandler.EditLastMessage)
	r.DELETE("/user/lastMessages/delete/:id", authHandler.CheckAuth, messageHandler.DeleteLastMessage)

	r.Run()
}
