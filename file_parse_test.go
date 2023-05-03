package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// language=JSON
var testPortsJson = `
{
  "DJJIB": {
    "name": "Djibouti",
    "city": "Djibouti",
    "country": "Djibouti",
    "alias": [],
    "regions": [],
    "coordinates": [
      43.1456475,
      11.5720765
    ],
    "province": "Djibouti",
    "timezone": "Africa/Djibouti",
    "unlocs": [
      "DJJIB",
      "DJPOD"
    ],
    "code": "77701"
  },
  "DJPOD": {
    "name": "Djibouti",
    "city": "Djibouti",
    "country": "Djibouti",
    "alias": [],
    "regions": [],
    "coordinates": [
      43.1456475,
      11.5720765
    ],
    "province": "Djibouti",
    "timezone": "Africa/Djibouti",
    "unlocs": [
      "DJPOD",
      "DJJIB"
    ],
    "code": "77701"
  },
  "ERR": {
    "name": "ERROR Unclosed quote
  }
}`

func Test_ParsePortsStream(t *testing.T) {
	ctx := context.Background()
	portsFile := strings.NewReader(testPortsJson)
	res := make(chan ParseStream, 10)

	go ParsePortsStream(ctx, portsFile, res)
	r1 := <-res
	assert.Equal(t, []string{"DJJIB", "DJPOD"}, r1.Value.Unlocs)
	r2 := <-res
	assert.Equal(t, []string{"DJPOD", "DJJIB"}, r2.Value.Unlocs)
	r3 := <-res
	assert.NotNil(t, r3.Error)
}

func Test_ParsePortsFile(t *testing.T) {
	ctx := context.Background()
	portsDb = NewInmemoryDb()
	totalProcessed, err := ParsePortsFile(ctx, "ports.json")
	assert.Nil(t, err)
	p1 := portsDb.FindPort(ctx, "DJJIB")
	assert.Equal(t, []string{"DJJIB", "DJPOD"}, p1.Unlocs)
	allPorts := portsDb.GetAll(ctx)
	assert.Equal(t, uint64(1632), totalProcessed)
	assert.Equal(t, 1624, len(allPorts))
}
