package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_ParsePortsStream(t *testing.T) {
	var portsFile io.Reader
	res := make(chan ParseStream, 4)

	go ParsePortsStream(portsFile, res)
	r1 := <-res
	r2 := <-res
	assert.Equal(t, "1", r1.Value.Unlocs[0])
	assert.NotNil(t, r2.Error)
}
