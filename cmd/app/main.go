package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/seldomhappy/sqlc-test/config"
	"github.com/seldomhappy/sqlc-test/internal/api/handler"
	"github.com/seldomhappy/sqlc-test/internal/infrastructure/database"
	"github.com/seldomhappy/sqlc-test/internal/infrastructure/repository"
	usecase "github.com/seldomhappy/sqlc-test/internal/usecase/author"
)

func main() {
	cfg := config.Load()

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Initialize infrastructure layer
	db := database.New(conn)
	authorRepo := repository.New(db.GetConn())

	// Initialize use cases
	listUC := usecase.NewListAuthorsUseCase(authorRepo)
	getUC := usecase.NewGetAuthorUseCase(authorRepo)
	createUC := usecase.NewCreateAuthorUseCase(authorRepo)
	updateUC := usecase.NewUpdateAuthorUseCase(authorRepo)
	deleteUC := usecase.NewDeleteAuthorUseCase(authorRepo)

	// Initialize HTTP handler and routes
	authorHandler := handler.NewAuthorHandler(listUC, getUC, createUC, updateUC, deleteUC)

	mux := http.NewServeMux()
	authorHandler.RegisterRoutes(mux)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start HTTP server
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
		Handler: mux,
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server on %s in %s mode", server.Addr, cfg.Environment)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
	log.Println("Server stopped gracefully")
}
