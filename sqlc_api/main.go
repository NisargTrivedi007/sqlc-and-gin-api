package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sqlc_api/api"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// jwt setup

	router := gin.Default()
	err := api.InitDB()
	if err != nil {
		panic(err)
	}
	defer api.CloseDB() // Make sure to create this function in your api package

	api.SetupRoutes(router)

	// Create a new HTTP server
	srv := &http.Server{
		Addr:    ":8080", // You can make this configurable
		Handler: router,
	}

	// Start server in a goroutine so it doesn't block shutdown handling
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to tell the server it has 5 seconds to finish
	// the requests it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
