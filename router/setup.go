package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/handlers"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

func Setup(db *gormdb.GormDB, jwtSecret string, asynqClient *asynq.Client, baseURL string) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	jwtSecretBytes := []byte(jwtSecret)
	audience := []string{baseURL}
	parseUserAuthToken := tokens.ParseUserAuth(jwtSecretBytes, audience, baseURL)
	checkAuth := middlewares.CheckAuth(parseUserAuthToken)
	parseEmailVerificationToken := tokens.ParseEmailVerification(jwtSecretBytes, baseURL, baseURL)
	createUserAuthToken := tokens.CreateUserAuth(jwtSecretBytes, audience, baseURL)

	// user stuff
	{
		user := r.Group("/user")
		user.POST("/register", handlers.RegisterUser(db, argon2id.CreateHash, parseEmailVerificationToken))
		user.POST("/verify-email", handlers.VerifyEmail(tasks.EnqueueTask(asynqClient), db))
		user.POST("/login", handlers.LoginUser(db, argon2id.ComparePasswordAndHash, createUserAuthToken))
		user.GET("/profile", checkAuth, handlers.GetUserData(db))
		user.PUT("/set-email-interval", checkAuth, handlers.SetEmailInterval(db))
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
		user.POST("/last-messages/add", checkAuth, handlers.AddLastMessage(db))
		user.GET("/last-messages", checkAuth, handlers.ListLastMessages(db))
		user.PATCH("/last-messages/edit/:id", checkAuth, handlers.EditLastMessage(db))
		user.DELETE("/last-messages/delete/:id", checkAuth, handlers.DeleteLastMessage(db))
	}
	return r
}
