package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/initializers"
	"github.com/gragorther/epigo/middlewares"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/scheduler"
	"github.com/gragorther/epigo/workers"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	dbconn, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = dbconn.AutoMigrate(&models.User{}, &models.LastMessage{}, &models.Group{}, &models.RecipientEmail{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	dbHandler := db.NewDBHandler(dbconn)
	redisAddr := os.Getenv("REDIS_ADDRESS")
	go workers.Run(redisAddr)
	go scheduler.Run(dbHandler, redisAddr)
	/*
		adminUsername := os.Getenv("ADMIN_USERNAME")
		adminPasswordHash, err := argon2id.CreateHash(os.Getenv("ADMIN_PASSWORD"), argon2id.DefaultParams)
		if err != nil {
			log.Fatalf("Failed to create admin account %v", err)
		}

		adminEmail := os.Getenv("ADMIN_EMAIL")
		res := dbconn.Create(&models.User{Username: adminUsername, Email: adminEmail, PasswordHash: adminPasswordHash, EmailCron: &adminCron, IsAdmin: true})
		if res.Error != nil {
			log.Print(res.Error)
		}
	*/
	r := gin.Default()
	userHandler := handlers.NewUserHandler(dbconn, dbHandler)
	authHandler := middlewares.NewAuthHandler(dbconn)
	groupHandler := handlers.NewGroupHandler(dbHandler)
	messageHandler := handlers.NewMessageHandler(dbHandler)

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
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	sqlDB, err := dbconn.DB()
	if err != nil {
		log.Fatal("failed to get sqldb ", err)
	}
	log.Print("closing db connection")
	sqlerr := sqlDB.Close()
	if sqlerr != nil {
		log.Printf("Failed to close db connection: %v", sqlerr)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
