package config

import (
	"encoding/json"
	"log"
	"os"
)

// PortsServiceConf is loaded from conf.json
// TODO The same config is used for import tool, needs for a split
type PortsServiceConf struct {
	DatabaseUrl   string
	ListenAddr    string
	Credentials   map[string]string
	PortsFilePath string
}

// LoadConfig Load settings from conf.json and optionally from conf.local.json (used for a local development)
func LoadConfig() *PortsServiceConf {
	conf := &PortsServiceConf{}
	confBytes, err := os.ReadFile("conf.json")
	if err != nil {
		log.Fatalf("Unable to read conf.json: %s\n", err)
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		log.Fatalf("Unable to read conf.json broken JSON: %s\n", err)
	}
	// override config fields from conf.local.json
	confBytes, err = os.ReadFile("conf.local.json")
	if err != nil {
		if os.IsNotExist(err) {
			return conf
		}
		log.Fatalf("Unable to read conf.local.json: %s\n", err)
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		log.Fatalf("Unable to read conf.json broken JSON: %s\n", err)
	}
	log.Printf("INFO config overriden from conf.local.json")
	return conf
}
