package main

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

	"github.com/http-pw-hash/api"
)

const (
	ServerName        = "Password Hash Server"
	ServerPort        = "8080"
	ShutdownTimeoutSc = 5
)

func main() {
	port := ServerPort

	if len(os.Args) > 1 {
		port = os.Args[1]
		if _, err := strconv.ParseInt(port, 10, 64); err != nil {
			log.Fatalf("ERROR: Bad port: %s, %s", port, err)
		}
	}

	hs := api.NewHashServer(syscall.Getpid())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: hs,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR: Unable to start server at port %s, %s", port, err)
		}
	}()
	log.Printf("[%d] %s Started, PORT: %s", hs.PID, ServerName, port)

	<-done
	log.Printf("[%d] %s Stopped, PORT: %s", hs.PID, ServerName, port)

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeoutSc*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("ERROR: [%d] %s Shutdown Failed:%+v", hs.PID, ServerName, err)
	}
	log.Printf("[%d] %s Shutdown Completed", hs.PID, ServerName)
}
