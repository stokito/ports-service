package main

import (
	"context"
	"log"
	"sync"
)

// InmemoryDb is an implementation of PortsDb that stores to an in-memory map
type InmemoryDb struct {
	sync.RWMutex
	// Map of unlocs to ports
	ports map[string]*Port
}

func NewInmemoryDb() *InmemoryDb {
	inmemDb := &InmemoryDb{
		ports: make(map[string]*Port, 2000),
	}
	return inmemDb
}

func (db *InmemoryDb) Connect(_ context.Context) error {
	log.Printf("INFO Connected to In-memory database\n")
	return nil
}

func (db *InmemoryDb) Close() {
	log.Printf("INFO In-memory database closed\n")
	db.ports = nil
}

// UpsertPort stores a Port to memory.
func (db *InmemoryDb) UpsertPort(_ context.Context, portUnloc string, port *Port) {
	db.Lock()
	defer db.Unlock()
	db.ports[portUnloc] = port
}

// FindPort retrieves a Port from memory by the portUnloc.
func (db *InmemoryDb) FindPort(_ context.Context, portUnloc string) *Port {
	db.RLock()
	defer db.RUnlock()
	p := db.ports[portUnloc]
	return p
}

func (db *InmemoryDb) GetAll(ctx context.Context) []*Port {
	allPorts := make([]*Port, 0, len(db.ports))
	for _, port := range db.ports {
		allPorts = append(allPorts, port)
	}
	return allPorts
}
