package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
)

func Setup(db gormdb.GormDBs) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	userHandler := handlers.NewUserHandler(userUserStore)
	authHandler := middlewares.NewAuthMiddleware(middlewareUserStore)
	groupHandler := handlers.NewGroupHandler(groupGroupStore, groupAuthStore)
	messageHandler := handlers.NewMessageHandler(messageMessageStore, messageAuthStore)

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
