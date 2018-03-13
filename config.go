package tdns

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ListenAddr     string // listen addr, eg. ":8080", "127.0.0.1:8080"
	DNSServer      string // dns server addr, eg. "8.8.8.8:53"
	IPDB           string // ip db file
	BadgerDir      string // dir to store badger data
	BadgerValueDir string // dir to stroe badger value log in
	LogDir         string // dir to save log
}

// LoadConfig read config file, unmarshal data to config struct
// LoadConfig returns a Config
func LoadConfig(file string) (*Config, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
