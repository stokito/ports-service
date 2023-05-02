package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func ParsePortsFile(portsFilePath string) error {
	portsFile, err := os.Open(portsFilePath)
	if err != nil {
		return err
	}
	defer portsFile.Close()

	var totalProcessed uint64

	// Communication
	res := make(chan ParseStream, 1024)

	go ParsePortsStream(portsFile, res)
	// Read from results stream...
	for got := range res {
		if got.Error != nil {
			if got.Error == io.EOF {
				break
			}
			log.Print(got.Error)
			continue
		}
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
	Value *Port
	Error error
}

func ParsePortsStream(portsFile io.Reader, res chan<- ParseStream) {
	defer close(res)
	port := &Port{
		Unlocs: []string{"1"},
	}
	res <- ParseStream{
		Value: port,
		Error: nil,
	}
	res <- ParseStream{
		Value: port,
		Error: errors.New("Error"),
	}
}
