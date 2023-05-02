package main

import (
	"context"
	"log"
)

func main() {
	log.Printf("Start Ports Service\n")
	ctx := context.Background()
	err := ParsePortsFile(ctx, "ports.json")
	if err != nil {
		log.Printf("Failed to parse ports file: %s\n", err)
		return
	}
}
