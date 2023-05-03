package pkg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InmemoryDB(t *testing.T) {
	ctx := context.Background()
	db := NewInmemoryDb()
	_ = db.Connect(ctx)
	defer db.Close()
	db.RemoveAll()
	defer db.RemoveAll()
	portToSave := &Port{
		Name:   "1",
		Unlocs: []string{"1", "2"},
	}
	portToSave2 := &Port{
		Name:   "2",
		Unlocs: []string{"2", "1"}, // unordered
	}
	db.UpsertPort(ctx, "1", portToSave)
	storedPort := db.FindPort(ctx, "1")
	assert.Equal(t, "1", storedPort.Name)
	assert.Equal(t, "1", storedPort.Unlocs[0])
	assert.Equal(t, "2", storedPort.Unlocs[1])
	allPorts := db.GetAll(ctx)
	assert.Equal(t, "1", allPorts[0].Name)
	assert.Equal(t, "1", allPorts[0].Unlocs[0])
	assert.Equal(t, "2", allPorts[0].Unlocs[1])
	assert.Equal(t, 1, len(allPorts))

	db.UpsertPort(ctx, "2", portToSave2)
	storedPort = db.FindPort(ctx, "2")
	assert.Equal(t, "2", storedPort.Name)
	assert.Equal(t, "1", storedPort.Unlocs[0])
	assert.Equal(t, "2", storedPort.Unlocs[1])
	allPorts = db.GetAll(ctx)
	assert.Equal(t, "2", allPorts[0].Name)
	assert.Equal(t, "1", allPorts[0].Unlocs[0])
	assert.Equal(t, "2", allPorts[0].Unlocs[1])
	assert.Equal(t, 1, len(allPorts))
}
