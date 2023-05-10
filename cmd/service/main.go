package main

import (
	"context"
	. "github.com/stokito/ports-service/internal/pkg/api"
	. "github.com/stokito/ports-service/internal/pkg/config"
	. "github.com/stokito/ports-service/internal/pkg/db"
	"log"
	"os"
	"os/signal"
)

// Starts a REST API service
// The service is intended for an internal use and uses a Basic Authorization.
// Only standard golang http server and mux is used for simplicity.
// Additional endpoints for profiling are exposed. They also requires the Basic Auth.
func main() {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancelFn()

	conf, err := LoadConfig()
	if err != nil {
		log.Fatalf("CRIT Unable to load config: %s\n", err)
	}
	err = InitDb(conf.DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to initialize DB: %s\n", err)
		return
	}
	dbErr := PortsDbConn.Connect(ctx)
	if dbErr != nil {
		log.Printf("CRIT Database connection failed: %s\n", dbErr)
		return
	}

	go func() {
		<-ctx.Done()
		StopApiServer()
	}()
	StartApiServer(conf.ListenAddr, conf.Credentials)
}
