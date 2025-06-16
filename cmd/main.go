package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.gllm.dev/vanity-go/di"
)

func main() {
	ctx := context.Background()
	slog.Info("Starting vanity-go server")

	server, err := di.ProvideRestServer()
	if err != nil {
		slog.Error("Failed to initialize dependencies", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		<-sigCh
		slog.Info("Received shutdown signal")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30)
		defer cancel()

		if err := server.Stop(shutdownCtx); err != nil {
			slog.Error("Failed to shutdown server gracefully", slog.String("error", err.Error()))
		}

		wg.Done()
	}()

	if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("Server stopped")
}
