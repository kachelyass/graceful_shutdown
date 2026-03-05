package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/m/v2/internal/database"
	"example.com/m/v2/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := database.Connect("postgresql://user:pass@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := database.NewTaskStore(db)
	taskHandler := handlers.NewTaskHandler(store)

	r := chi.NewRouter()
	r.Mount("/tasks", taskHandler.Routes())

	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Println("server started on :8080")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	case <-shutdownCtx.Done():
		log.Println("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)

		if err := server.Close(); err != nil {
			log.Printf("forced server close failed: %v", err)
		}
	}

	log.Println("server stopped")
}
