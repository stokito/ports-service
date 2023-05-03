package main

import (
	"context"
	"log"
)

type PortsServiceConf struct {
	DatabaseUrl string
}

var conf *PortsServiceConf

func main() {
	conf = &PortsServiceConf{DatabaseUrl: "postgres://postgres:postgres@127.0.0.1:5432/portsdb?search_path=ports_schema"}
	log.Printf("Start Ports Service\n")
	ctx := context.Background()
	dbErr := dbConnect(ctx)
	if dbErr != nil {
		log.Printf("CRIT Database connection failed: %s\n", dbErr)
		return
	}
	defer dbClose()
	err := ParsePortsFile(ctx, "ports.json")
	if err != nil {
		log.Printf("CRIT Failed to parse ports file: %s\n", err)
		return
	}
}
