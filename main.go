package main

import "log"

func main() {
	log.Printf("Start Ports Service\n")
	err := ParsePortsFile("ports.json")
	if err != nil {
		log.Printf("Failed to parse ports file: %s\n", err)
		return
	}
}
