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

var ErrorBadUnloc = errors.New("bad Unloc")

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
		portsDb.UpsertPort(ctx, got.Unloc, got.Value)
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
