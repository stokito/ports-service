package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInmemoryDB_Connect(t *testing.T) {
	ctx := context.Background()
	db := NewInmemoryDb()
	_ = db.Connect(ctx)
	defer db.Close()
	portToSave := &Port{
		Unlocs: []string{"1", "2"},
	}
	db.UpsertPort(ctx, "1", portToSave)
	storedPort := db.FindPort(ctx, "1")
	assert.Equal(t, "1", storedPort.Unlocs[0])
	allPorts := db.GetAll(ctx)
	assert.Equal(t, "1", allPorts[0].Unlocs[0])
	assert.Equal(t, "2", allPorts[0].Unlocs[1])
	assert.Equal(t, 1, len(allPorts))

	db.UpsertPort(ctx, "2", portToSave)
	storedPort = db.FindPort(ctx, "2")
	assert.Equal(t, "1", storedPort.Unlocs[0])
	assert.Equal(t, "2", allPorts[0].Unlocs[1])
	allPorts = db.GetAll(ctx)
	assert.Equal(t, 2, len(allPorts))
}
