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

	"github.com/godbus/dbus/v5"
	"github.com/nickrobison/terraform-linux-provider/server/firewall"
	"github.com/nickrobison/terraform-linux-provider/server/middleware"
	"github.com/nickrobison/terraform-linux-provider/server/zfs"
	"github.com/rs/zerolog"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	middleware.SetupLogging(w, zerolog.DebugLevel)
	log := middleware.Logger()

	log.Print("Hello world!")

	// Connect to DBUS
	// Should be session bus
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to bus")
		return err
	}

	zfsClient, err := zfs.NewZfsClient(conn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to dbus object")
		return err
	}
	zfsVersion, err := zfsClient.Version()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Initialized Zfs client with version %s", zfsVersion)

	firewallClient, err := firewall.NewFirewallClient(conn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to firewall dbus object")
		return err
	}
	firewallVersion, err := firewallClient.Version()
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Initialized Firewall client with version %s", firewallVersion)

	srv := newServer(zfsClient, firewallClient)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: srv,
	}

	go func() {
		log.Info().Msgf("Listenting on %s", httpServer.Addr)
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
