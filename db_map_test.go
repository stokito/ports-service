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
	db.Close(ctx)
	portToSave := &Port{
		Unlocs: []string{"1"},
	}
	db.UpsertPort(ctx, "1", portToSave)
	storedPort := db.FindPort(ctx, "1")
	assert.Equal(t, "1", storedPort.Unlocs[0])
}
