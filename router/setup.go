package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/middlewares"
)

func Setup(db *gormdb.GormDB, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	checkAuth := middlewares.CheckAuth(db, jwtSecret)

	// user stuff
	r.POST("/user/register", handlers.RegisterUser(db, argon2id.CreateHash))
	r.POST("/user/login", handlers.LoginUser(db, argon2id.ComparePasswordAndHash, jwtSecret))
	r.GET("/user/profile", checkAuth, handlers.GetUserData())
	r.PUT("/user/setEmailInterval", checkAuth, handlers.SetEmailInterval(db))

	// groups
	r.DELETE("/user/groups/delete/:id", checkAuth, handlers.DeleteGroup(db))
	r.POST("/user/groups/add", checkAuth, handlers.AddGroup(db))
	r.GET("/user/groups", checkAuth, handlers.ListGroups(db)) // list groups
	r.PATCH("/user/groups/edit/:id", checkAuth, handlers.EditGroup(db))

	// lastMessages
	r.POST("/user/lastMessages/add", checkAuth, handlers.AddLastMessage(db))
	r.GET("/user/lastMessages", checkAuth, handlers.ListLastMessages(db))
	r.PATCH("/user/lastMessages/edit/:id", checkAuth, handlers.EditLastMessage(db))
	r.DELETE("/user/lastMessages/delete/:id", checkAuth, handlers.DeleteLastMessage(db))
	return r
}
