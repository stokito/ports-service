package db

import (
	"context"
	"errors"
	. "github.com/stokito/ports-service/internal/pkg/domain"
	"strings"
)

// PortsDb is a generic DB interface
type PortsDb interface {
	Connect(ctx context.Context) error
	Close()
	// UpsertPort insert or update a Port. The portUnloc may be empty then port.Unloc will be used
	UpsertPort(ctx context.Context, portUnloc string, port *Port)
	// FindPort retrieves a Port by the portUnloc.
	FindPort(ctx context.Context, portUnloc string) *Port
	// GetAll Return a list of all stored ports. Their order is not guaranteed
	GetAll(ctx context.Context) []*Port
	RemoveAll()
}

var PortsDbConn PortsDb

func InitDb(dbUrl string) error {
	if dbUrl == "" {
		return errors.New("database is not configured")
	}
	if strings.HasPrefix(dbUrl, "mem://") {
		PortsDbConn = NewInmemoryDb()
		return nil
	} else if strings.HasPrefix(dbUrl, "postgres://") {
		PortsDbConn = NewPostgresDb(dbUrl)
		return nil
	}
	return errors.New("database is not supported")
}
