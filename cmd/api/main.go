package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/handler"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func main() {
	if err := run(os.Getenv, http.ListenAndServe); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func run(getenv func(string) string, listenAndServe func(string, http.Handler) error) error {
	port := getEnvWith(getenv, "PORT", "8080")
	r := mux.NewRouter()

	// Use PostgreSQL if DB_HOST is set, otherwise fall back to in-memory store
	dbHost := getenv("DB_HOST")
	if dbHost != "" {
		pgStore, err := store.NewPostgresStore(
			dbHost,
			getEnvWith(getenv, "DB_PORT", "5432"),
			getEnvWith(getenv, "DB_USER", "catalog"),
			getEnvWith(getenv, "DB_PASSWORD", "catalog123"),
			getEnvWith(getenv, "DB_NAME", "productcatalog"),
		)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		if err := pgStore.EnsureTable(); err != nil {
			_ = pgStore.DB.Close()
			return fmt.Errorf("failed to create table: %w", err)
		}
		defer func() {
			_ = pgStore.DB.Close()
		}()

		h := handler.NewPostgresHandler(pgStore)
		h.RegisterRoutes(r)
		fmt.Printf("Product Catalog API (PostgreSQL) listening on :%s\n", port)
	} else {
		memStore := store.NewMemoryStore()
		h := handler.NewHandler(memStore)
		h.RegisterRoutes(r)
		fmt.Printf("Product Catalog API (in-memory) listening on :%s\n", port)
	}

	return listenAndServe(net.JoinHostPort("", port), r)
}

func getEnv(key, fallback string) string {
	return getEnvWith(os.Getenv, key, fallback)
}

func getEnvWith(getenv func(string) string, key, fallback string) string {
	if v := getenv(key); v != "" {
		return v
	}
	return fallback
}
