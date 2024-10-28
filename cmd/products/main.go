package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func handleGetUserProducts(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	log.Printf("Getting products for user %s...", userID)

	// get user products from DB
	time.Sleep(4 * time.Second)
	products := []string{"product1", "product2"}

	writeJSON(w, http.StatusOK, products)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func main() {
	ctx := context.Background()

	wg, wgCtx := errgroup.WithContext(ctx)

	stopCtx, _ := signal.NotifyContext(wgCtx, syscall.SIGTERM, syscall.SIGINT)

	router := http.NewServeMux()
	router.HandleFunc("GET /users/{id}/products", handleGetUserProducts)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// Start the server.
	wg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Println("server stopped")
		return nil
	})

	// Gracefully shut down the server.
	wg.Go(func() error {
		<-stopCtx.Done()
		log.Println("shutting down server")
		server.Shutdown(ctx)
		return nil
	})

	log.Println("server is running on", server.Addr)

	// Wait for all go routines to finish.
	if err := wg.Wait(); err != nil {
		log.Printf("stopped with error: %v\n", err)
		os.Exit(1)
	}
}
