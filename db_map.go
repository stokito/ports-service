package main

import (
	"context"
	"log"
	"sync"
)

// InmemoryDb is an implementation of PortsDb that stores to an in-memory map
type InmemoryDb struct {
	sync.RWMutex
	// Map of unlocs to ports. All unlocs of the same port will point to the same port
	ports map[string]*Port
	// List of all unique ports that were added. Order is not guaranteed but for testing may be assumed same as were added
	portsList []*Port
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
	oldPort := db.ports[portUnloc]
	// replace a reference for all unlocs of the old port to the new one
	if oldPort != nil {
		// first step is to clear all existing unlocs of old port if any
		// after this the oldPort instance wouldn't have pointer to it and should be freed by GC
		for _, oldUnloc := range oldPort.Unlocs {
			delete(db.ports, oldUnloc)
		}
	} else {
		// append the port to ports list
		db.portsList = append(db.portsList, port)
	}
	// add all port's unlocs
	for _, newUnloc := range port.Unlocs {
		db.ports[newUnloc] = port
	}
}

// FindPort retrieves a Port from memory by the portUnloc.
func (db *InmemoryDb) FindPort(_ context.Context, portUnloc string) *Port {
	db.RLock()
	defer db.RUnlock()
	p := db.ports[portUnloc]
	return p
}

func (db *InmemoryDb) GetAll(ctx context.Context) []*Port {
	return db.portsList
}
