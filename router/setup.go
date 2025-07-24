package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
)

func Setup(dbHandler *db.DBHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	userHandler := handlers.NewUserHandler(dbHandler.Users)
	authHandler := middlewares.NewAuthMiddleware(dbHandler.Users)
	groupHandler := handlers.NewGroupHandler(dbHandler.Groups, dbHandler.Auth)
	messageHandler := handlers.NewMessageHandler(dbHandler.Messages, dbHandler.Auth)

	// user stuff
	r.POST("/user/register", userHandler.RegisterUser)
	r.POST("/user/login", userHandler.LoginUser)
	r.GET("/user/profile", authHandler.CheckAuth, userHandler.GetUserProfile)
	r.PUT("/user/setEmailInterval", authHandler.CheckAuth, userHandler.SetEmailInterval)

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
	return r
}
