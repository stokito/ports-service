package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

var ErrorBadUnloc = errors.New("bad Unloc")

func ParsePortsFile(ctx context.Context, portsFilePath string) (uint64, error) {
	portsFile, err := os.Open(portsFilePath)
	if err != nil {
		return 0, err
	}
	defer portsFile.Close()

	var totalProcessed uint64

	// Communication
	res := make(chan ParseStream, 1024)

	go ParsePortsStream(ctx, portsFile, res)
	// Read from results stream...
processing:
	for {
		select {
		case <-ctx.Done():
			break processing
		case got, ok := <-res:
			if !ok {
				break processing
			}
			if got.Error != nil {
				log.Printf("ERR on process: %s\n", got.Error)
				continue
			}
			if got.Value == nil {
				continue
			}
			portsDb.UpsertPort(ctx, got.Unloc, got.Value)
			totalProcessed++
		}
	}
	return totalProcessed, nil
}

// A ParseStream is used to stream back results.
// Either Error or Value will be set on returned results.
type ParseStream struct {
	// Unloc code, can be empty on parsing error
	Unloc string
	Value *Port
	Error error
}

func ParsePortsStream(ctx context.Context, portsFile io.Reader, res chan<- ParseStream) {
	defer close(res)

	jsonDec := json.NewDecoder(portsFile)
	_, _ = jsonDec.Token() // skip first {
	for jsonDec.More() {
		unlocToken, err := jsonDec.Token()
		if err != nil {
			res <- ParseStream{Error: err}
			continue
		}
		unloc, ok := unlocToken.(string)
		if !ok {
			res <- ParseStream{Error: ErrorBadUnloc}
			continue
		}
		//TODO Get the raw JSON of the Port instead of mapping to the struct.
		// The port structure is not fully specified and can contain many other fields.
		// The Golang's decoder doesn't support getting of the raw json in effective way.
		port := &Port{}
		err = jsonDec.Decode(port)
		if err != nil {
			res <- ParseStream{Unloc: unloc, Error: err}
			continue
		}
		res <- ParseStream{Unloc: unloc, Value: port}
	}
}
