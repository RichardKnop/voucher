package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RichardKnop/voucher/server"
)

// TODO - get from environment
var (
	wait = time.Duration(5 * time.Second) // for graceful shutdown

	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: time.Duration(5 * time.Second),
	}
)

// RunServer runs an HTTP server
func RunServer() error {
	appServer := server.New()

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
