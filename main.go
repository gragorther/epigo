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
	"github.com/gragorther/epigo/asynq/scheduler"
	"github.com/gragorther/epigo/asynq/workers"
	"github.com/gragorther/epigo/config"
	"github.com/gragorther/epigo/database/gormdb"
	"github.com/gragorther/epigo/database/initializers"
	"github.com/gragorther/epigo/logger"
	"github.com/gragorther/epigo/router"
)

func main() {
	config, err := config.Get()
	if err != nil {
		os.Exit(1)
	}

	_ = logger.Configure(config.Production, os.Stdout)

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

	dbHandler := gormdb.NewGormDB(dbconn)

	go workers.Run(config.RedisAddress)
	go scheduler.Run(dbHandler, config.RedisAddress)

	r := router.Setup(dbHandler, config.JWTSecret)

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
