package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Print("Hello world!")

	srv := newServer()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: srv,
	}

	go func() {
		log.Printf("Listenting on %s\n", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to shutdown server")
		}
	}()
	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	err := run(ctx, os.Stdout, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
