package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/RichardKnop/voucher/server"
	"github.com/RichardKnop/voucher/service"
	"github.com/go-redis/redis"
)

var (
	redisHost = os.Getenv("REDIS_HOST")
	redisDB   = os.Getenv("REDIS_DB")

	wait = time.Duration(5 * time.Second) // for graceful shutdown
)

// RunServer runs an HTTP server
func RunServer() error {
	// Connect to redis
	if redisHost == "" {
		log.Println("REDIS_HOST environment variable empty, using default")
		redisHost = "localhost:6379"
	}
	if redisDB == "" {
		log.Println("REDIS_DB environment variable empty, using default")
		redisDB = "localhost:6379"
	}
	redisDBInt, _ := strconv.Atoi(redisDB)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost, // use default Addr
		Password: "",        // no password set
		DB:       redisDBInt,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong)

	// Create service
	service := service.New(redisClient)

	// Create app server
	appServer := server.New(service)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      appServer,
	}

	// Run the server in a goroutine so that it doesn't block
	shutdown := make(chan int)

	go func() {
		log.Println("Running server at 0.0.0.0:8080")

		srv.ListenAndServe()

		shutdown <- 1
	}()

	sigChan := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (docker stop)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-sigChan
		log.Printf("Signal received: %v", s)

		log.Println("Shutting down...")

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		srv.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
	}()

	<-shutdown

	return nil
}
