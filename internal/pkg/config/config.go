package config

import (
	"encoding/json"
	"github.com/pkg/errors"
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
func LoadConfig() (*PortsServiceConf, error) {
	conf := &PortsServiceConf{}
	confBytes, err := os.ReadFile("conf.json")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read conf.json")
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read conf.json broken JSON")
	}
	// override config fields from conf.local.json
	confBytes, err = os.ReadFile("conf.local.json")
	if err != nil {
		if os.IsNotExist(err) {
			return conf, nil
		}
		return nil, errors.Wrap(err, "Unable to read conf.local.json")
	}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read conf.json broken JSON")
	}
	log.Printf("INFO config overriden from conf.local.json")
	return conf, err
}
