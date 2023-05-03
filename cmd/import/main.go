package main

import (
	"context"
	. "github.com/stokito/ports-service/internal/pkg/config"
	. "github.com/stokito/ports-service/internal/pkg/db"
	. "github.com/stokito/ports-service/internal/pkg/parser"
	"log"
)

func main() {
	conf := LoadConfig()
	log.Printf("INFO Start Ports Import\n")
	ctx := context.Background()
	dbInitErr := InitDb(conf.DatabaseUrl)
	if dbInitErr != nil {
		log.Printf("CRIT Database configuration failed: %s\n", dbInitErr)
		return
	}
	dbErr := PortsDbConn.Connect(ctx)
	if dbErr != nil {
		log.Printf("CRIT Database connection failed: %s\n", dbErr)
		return
	}
	defer PortsDbConn.Close()
	log.Printf("INFO Importing %s\n", conf.PortsFilePath)
	totalProcessed, err := ParsePortsFile(ctx, conf.PortsFilePath)
	if err != nil {
		log.Printf("CRIT Failed to parse ports file: %s\n", err)
		return
	}
	log.Printf("INFO Total Processed %d\n", totalProcessed)
}
