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

	"github.com/Mi7teR/aggregator/internal/task/delivery/api"
	"github.com/Mi7teR/aggregator/internal/task/repository"
	"github.com/Mi7teR/aggregator/internal/task/service"
)

func main() {
	timeoutENV := os.Getenv("TIMEOUT")
	httpPortENV := os.Getenv("HTTP_PORT")

	timeout, err := time.ParseDuration(timeoutENV)
	if err != nil {
		log.Fatalln(fmt.Errorf("cant parse timeout %w", err))
	}

	repo := repository.NewTaskInMemoryRepository()
	s := service.NewService(repo, timeout)
	handler := api.NewHandler(s)
	r := api.NewRouter(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPortENV),
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
	log.Print("Server Exited Properly")
}
