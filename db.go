package main

import (
	"context"
	"errors"
	"strings"
)

// PortsDb is a generic DB interface
type PortsDb interface {
	Connect(ctx context.Context) error
	Close()
	UpsertPort(ctx context.Context, portUnloc string, port *Port)
	// FindPort retrieves a Port by the portUnloc.
	FindPort(ctx context.Context, portUnloc string) *Port
	// GetAll Return a list of all stored ports. Their order is not guaranteed
	GetAll(ctx context.Context) []*Port
}

var portsDb PortsDb

func InitDb() error {
	if conf.DatabaseUrl == "" {
		return errors.New("database is not configured")
	}
	if strings.HasPrefix(conf.DatabaseUrl, "mem://") {
		portsDb = NewInmemoryDb()
		return nil
	} else if strings.HasPrefix(conf.DatabaseUrl, "postgres://") {
		portsDb = NewPostgresDb(conf.DatabaseUrl)
		return nil
	}
	return errors.New("database is not supported")
}
