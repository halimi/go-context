package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type User struct {
	ID       string
	Name     string
	Products []string
}

type authenticatedUserKeyType struct{}

var authenticatedUserKey = authenticatedUserKeyType{}

// MiddlewareFunc is a function that wraps an HTTP handler.
type MiddlewareFunc func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, middleware ...MiddlewareFunc) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// WithAuthenticatedUser returns a new request with the authenticated user set.
func WithAuthenticatedUser(r *http.Request, username string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, username))
}

// AuthenticatedUser returns the authenticated user from a request context.
func AuthenticatedUser(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(authenticatedUserKey).(string)
	return username, ok
}

func Authenticate() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Shop"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			r = WithAuthenticatedUser(r, username)
			next.ServeHTTP(w, r)
		})
	}
}

func handleGetUserProducts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rCtx := r.Context()
		userID := r.PathValue("id")

		userName, ok := AuthenticatedUser(rCtx)
		if !ok {
			http.Error(w, "could not get username from context", http.StatusInternalServerError)
			return
		}

		user := User{
			ID:       userID,
			Name:     userName,
			Products: []string{},
		}

		ctx, cancel := context.WithTimeout(rCtx, 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:8081/users/%s/products", userID), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(data, &user.Products); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, user)
	})
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
	router.Handle("GET /users/{id}", applyMiddleware(
		handleGetUserProducts(),
		Authenticate(),
	))

	server := http.Server{
		Addr:    ":8080",
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
