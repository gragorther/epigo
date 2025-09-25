package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/asynq/tasks"
	"github.com/gragorther/epigo/database/db"
	"github.com/gragorther/epigo/handlers/groups"
	"github.com/gragorther/epigo/handlers/messages"
	"github.com/gragorther/epigo/handlers/users"
	argon2id "github.com/gragorther/epigo/hash"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/tokens"
)

func Setup(db interface {
	CheckIfUserExistsByUsername(ctx context.Context, username string) (bool, error)
	CheckIfUserExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateUserInterval(context.Context, uint, string) error
	UserByID(context.Context, uint) (db.User, error)
	UserIDAndPasswordHashByUsername(ctx context.Context, username string) (user db.UserIDAndPasswordHash, err error)
	CreateUser(context.Context, db.CreateUserInput) error
	UserAuthorizationForLastMessage(ctx context.Context, messageID uint, userID uint) (bool, error)
	DeleteLastMessageByID(ctx context.Context, id uint) error
	CanUserEditLastmessage(ctx context.Context, userID uint, messageID uint, groupIDs []uint) (authorized bool, err error)
	UpdateLastMessage(ctx context.Context, id uint, m db.UpdateLastMessage) error
	LastMessagesByUserID(ctx context.Context, userID uint) (lastMessages []db.LastMessage, err error)
	CreateLastMessage(ctx context.Context, message db.CreateLastMessage) error
	CanUserEditGroup(ctx context.Context, userID uint, groupID uint, lastMessageIDs []uint) (authorized bool, err error)
	UpdateGroup(ctx context.Context, id uint, group db.UpdateGroup) error
	GroupsByUserID(ctx context.Context, userID uint) (groups []db.Group, err error)
	DeleteGroupByID(ctx context.Context, id uint) error
	UserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error)
	CreateGroup(context.Context, db.CreateGroup) error
	CheckIfUserExistsByUsernameAndEmail(ctx context.Context, username string, email string) (bool, error)
}, queue interface {
	UpdateLastMessage(id uint, m db.UpdateLastMessage) error
	SendVerificationEmail(email string) error
	UpdateUserInterval(uint, string) error
	CreateGroup(db.CreateGroup) error
	DeleteGroupByID(id uint) error
	UpdateGroup(id uint, group db.UpdateGroup) error
	CreateLastMessage(message db.CreateLastMessage) error
	DeleteLastMessageByID(id uint) error
	CreateUser(db.CreateUserInput) error
}, jwtSecret string, enqueueTask tasks.TaskEnqueueFunc, baseURL string, minDurationBetweenEmail time.Duration,
) *gin.Engine {
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
		user.POST("/register", users.Register(db, queue, argon2id.CreateHash, parseEmailVerificationToken))
		user.POST("/verify-email", users.VerifyEmail(queue, db))
		user.POST("/login", users.Login(db, argon2id.ComparePasswordAndHash, createUserAuthToken))
		user.GET("/profile", checkAuth, users.GetData(db))
		user.PUT("/set-email-interval", checkAuth, users.SetEmailInterval(queue, minDurationBetweenEmail))

		// groups
		user.DELETE("/groups/:id", checkAuth, groups.Delete(db, queue))
		user.POST("/groups", checkAuth, groups.Add(queue))
		user.GET("/groups", checkAuth, groups.List(db)) // list groups
		user.PATCH("/groups/:id", checkAuth, groups.Edit(db, queue))
		user.GET("/groups/:id", checkAuth)

		// lastMessages
		user.POST("/last-messages", checkAuth, messages.Add(db, queue))
		user.GET("/last-messages", checkAuth, messages.List(db))
		user.PATCH("/last-messages/:id", checkAuth, messages.Edit(db, queue))
		user.DELETE("/last-messages/:id", checkAuth, messages.Delete(db, queue))
	}
	return r
}
