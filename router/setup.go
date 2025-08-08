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
	{
		user := r.Group("/user")
		user.POST("/register", handlers.RegisterUser(db, argon2id.CreateHash))
		user.POST("/login", handlers.LoginUser(db, argon2id.ComparePasswordAndHash, jwtSecret))
		user.GET("/profile", checkAuth, handlers.GetUserData())
		user.PUT("/setEmailInterval", checkAuth, handlers.SetEmailInterval(db))
		{
			profile := user.Group("/profile")
			profile.POST("/create", checkAuth, handlers.CreateProfile(db))
			profile.PATCH("/edit", checkAuth, handlers.UpdateProfile(db))
		}
		// groups
		user.DELETE("/groups/delete/:id", checkAuth, handlers.DeleteGroup(db))
		user.POST("/groups/add", checkAuth, handlers.AddGroup(db))
		user.GET("/groups", checkAuth, handlers.ListGroups(db)) // list groups
		user.PATCH("/groups/edit/:id", checkAuth, handlers.EditGroup(db))

		// lastMessages
		user.POST("/lastMessages/add", checkAuth, handlers.AddLastMessage(db))
		user.GET("/lastMessages", checkAuth, handlers.ListLastMessages(db))
		user.PATCH("/lastMessages/edit/:id", checkAuth, handlers.EditLastMessage(db))
		user.DELETE("/lastMessages/delete/:id", checkAuth, handlers.DeleteLastMessage(db))
	}
	return r
}
