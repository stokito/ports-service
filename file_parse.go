package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var ErrorBadUnlock = errors.New("bad Unlock")

func ParsePortsFile(ctx context.Context, portsFilePath string) error {
	portsFile, err := os.Open(portsFilePath)
	if err != nil {
		return err
	}
	defer portsFile.Close()

	var totalProcessed uint64

	// Communication
	res := make(chan ParseStream, 1024)

	go ParsePortsStream(ctx, portsFile, res)
	// Read from results stream...
	for got := range res {
		select {
		case <-ctx.Done():
			break
		}
		if got.Error != nil {
			if got.Error == io.EOF {
				break
			}
			log.Print(got.Error)
			continue
		}
		portsDb.UpsertPort(ctx, got.Unlock, got.Value)
		totalProcessed++
	}
	fmt.Println("Total Processed ", totalProcessed)
	return nil
}

type Port struct {
	Unlocs []string `json:"unlocs"`
}

// A ParseStream is used to stream back results.
// Either Error or Value will be set on returned results.
type ParseStream struct {
	// Unlock code, can be empty on parsing error
	Unlock string
	Value  *Port
	Error  error
}

func ParsePortsStream(ctx context.Context, portsFile io.Reader, res chan<- ParseStream) {
	defer close(res)

	jsonDec := json.NewDecoder(portsFile)
	_, _ = jsonDec.Token() // skip first {
	for jsonDec.More() {
		unlockToken, err := jsonDec.Token()
		if err != nil {
			res <- ParseStream{Error: err}
			continue
		}
		unlock, ok := unlockToken.(string)
		if !ok {
			res <- ParseStream{Error: ErrorBadUnlock}
			continue
		}
		//TODO Get the raw JSON of the Port instead of mapping to the struct.
		// The port structure is not fully specified and can contain many other fields.
		// The Golang's decoder doesn't support getting of the raw json in effective way.
		port := &Port{}
		err = jsonDec.Decode(port)
		if err != nil {
			res <- ParseStream{Unlock: unlock, Error: err}
			continue
		}
		res <- ParseStream{Unlock: unlock, Value: port}
	}
}
