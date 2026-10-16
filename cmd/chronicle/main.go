package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chronicle/internal/api"
	"chronicle/internal/research"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	store, err := openWorkspace(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	if err := prepareWorkspace(ctx, store); err != nil {
		log.Fatal(err)
	}
	manager := research.NewManager(store)
	mux := http.NewServeMux()
	mux.Handle("/api/", api.New(store, manager))
	mux.Handle("/", http.FileServer(http.Dir("web/build")))
	address := os.Getenv("CHRONICLE_ADDR")
	if address == "" {
		address = "127.0.0.1:8080"
	}
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(timeout)
	}()
	log.Printf("Chronicle is available at http://%s", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}
	manager.Shutdown(context.Background())
}
