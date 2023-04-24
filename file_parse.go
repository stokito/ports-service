package main

import "os"

func ParsePortsFile(portsFilePath string) error {
	portsFile, err := os.Open(portsFilePath)
	if err != nil {
		return err
	}
	defer portsFile.Close()
	return nil
}
