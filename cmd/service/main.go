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

func main() {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancelFn()

	conf := LoadConfig()
	err := InitDb(conf.DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to initialize DB\n")
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
