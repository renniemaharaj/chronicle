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
	"chronicle/internal/config"
	"chronicle/internal/research"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] != "healthcheck" || len(os.Args) != 2 {
			log.Fatal("usage: chronicle [healthcheck]")
		}

		if err := healthcheck(settings.Address); err != nil {
			log.Fatal(err)
		}

		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store, err := openWorkspace(ctx, settings)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if err := prepareWorkspace(ctx, store); err != nil {
		log.Fatal(err)
	}

	manager := research.NewManager(store)
	mux := http.NewServeMux()
	mux.Handle("/api/", api.New(store, manager, settings))
	mux.Handle("/", http.FileServer(http.Dir(settings.StaticDir)))

	server := &http.Server{
		Addr:              settings.Address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		<-ctx.Done()

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.Shutdown(timeout)
	}()
	log.Printf("Chronicle is available at http://%s", settings.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}

	manager.Shutdown(context.Background())
}
