package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/middlewares"
)

func Setup(db *gormdb.GormDB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	// user stuff
	r.POST("/user/register", handlers.RegisterUser(db))
	r.POST("/user/login", handlers.LoginUser(db))
	r.GET("/user/profile", middlewares.CheckAuth(db), handlers.GetUserProfile())
	r.PUT("/user/setEmailInterval", middlewares.CheckAuth(db), handlers.SetEmailInterval(db))

	// groups
	r.DELETE("/user/groups/delete/:id", middlewares.CheckAuth(db), handlers.DeleteGroup(db))
	r.POST("/user/groups/add", middlewares.CheckAuth(db), handlers.AddGroup(db))
	r.GET("/user/groups", middlewares.CheckAuth(db), handlers.ListGroups(db)) // list groups
	r.PATCH("/user/groups/edit/:id", middlewares.CheckAuth(db), handlers.EditGroup(db))

	// lastMessages
	r.POST("/user/lastMessages/add", middlewares.CheckAuth(db), handlers.AddLastMessage(db))
	r.GET("/user/lastMessages", middlewares.CheckAuth(db), handlers.ListLastMessages(db))
	r.PATCH("/user/lastMessages/edit/:id", middlewares.CheckAuth(db), handlers.EditLastMessage(db))
	r.DELETE("/user/lastMessages/delete/:id", middlewares.CheckAuth(db), handlers.DeleteLastMessage(db))
	return r
}
