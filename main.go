package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/asynq/scheduler"
	"github.com/gragorther/epigo/asynq/workers"
	"github.com/gragorther/epigo/config"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/database/initializers"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/logger"
	"github.com/gragorther/epigo/router"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

func main() {
	config, err := config.Get()
	jwtSecret := []byte(config.JWTSecret)
	if err != nil {
		os.Exit(1)
	}

	_ = logger.Configure(config.Production, os.Stdout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if config.GinMode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	dbconn, err := initializers.ConnectDB(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	err = initializers.Migrate(dbconn)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	dbHandler := gormdb.NewGormDB(dbconn)
	emailClient, err := email.NewClient(config.Email.Host, config.Email.Port, config.Email.Password, config.Email.Username)
	defer func() {
		if err := emailClient.Close(); err != nil {
			log.Fatalf("failed to close email client: %v", err)
		}
	}()
	if err != nil {
		log.Fatalf("failed to run email client: %v", err)
	}
	createEmailVerificationToken := tokens.CreateEmailVerification(jwtSecret, config.BaseURL, config.BaseURL)
	createUserLifeStatusToken := tokens.CreateUserLifeStatus(jwtSecret, []string{config.BaseURL}, config.BaseURL)
	emailService := email.NewEmailService(emailClient, config.Email.From)
	redisClientOpt := asynq.RedisClientOpt{Addr: config.Redis.Address, Username: config.Redis.Address, Password: config.Redis.Password, DB: config.Redis.DB}
	go workers.Run(ctx, redisClientOpt, jwtSecret, emailService, fmt.Sprintf("%v/user/register", config.BaseURL), createEmailVerificationToken, createUserLifeStatusToken, fmt.Sprintf("%s/user/life/verify", config.BaseURL))
	go scheduler.Run(dbHandler, redisClientOpt)
	asynqClient := asynq.NewClient(redisClientOpt)

	r := router.Setup(dbHandler, config.JWTSecret, asynqClient, config.BaseURL, config.MinDurationBetweenEmails)

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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
