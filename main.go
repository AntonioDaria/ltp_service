package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/antoniodaria/ltp_service/api"
	"github.com/antoniodaria/ltp_service/clients/kraken"
)

func main() {
	// HTTP Server
	ctx := context.Background()

	krakenClient := kraken.NewClient(http.DefaultClient, "https://api.kraken.com")
	// Start the server
	server := api.NewServer(krakenClient)

	// Graceful shutdown
	/* Listen for SIGINT, SIGTERM, SIGQUIT
	This will gracefully shutdown the server without interrupting
	any active connections or requests being processed by the server.
	This is a good practice to ensure that the server is not abruptly terminated. */
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-osSignals
		fmt.Fprint(os.Stdout, "🛑 Shutting down Server")
		server.Shutdown()
		fmt.Fprint(os.Stdout, "✅ HTTP Server shutdown complete")
	}()

	fmt.Println("🚀 Starting HTTP Server")
	server.StartAndListen(ctx)
}
