package main

import (
	"context"
	"encoding/json"
	. "github.com/stokito/ports-service/internal/pkg"
	"log"
	"os"
)

type PortsServiceConf struct {
	DatabaseUrl   string
	PortsFilePath string
}

var conf *PortsServiceConf

func main() {
	loadConfig()
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

// loadConfig Load settings from conf.json and optionally from conf.local.json (used for a local development)
func loadConfig() {
	conf = &PortsServiceConf{}
	confBytes, err := os.ReadFile("conf.json")
	if err != nil {
		log.Fatalf("Unable to read conf.json: %s\n", err)
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		log.Fatalf("Unable to read conf.json broken JSON: %s\n", err)
	}
	// override config fields from conf.local.json
	confBytes, err = os.ReadFile("conf.local.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatalf("Unable to read conf.local.json: %s\n", err)
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		log.Fatalf("Unable to read conf.json broken JSON: %s\n", err)
	}
	log.Printf("INFO config overriden from conf.local.json")
}