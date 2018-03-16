package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config app config
type Config struct {
	RPCAddr        string // listen addr, eg. ":8080", "127.0.0.1:8080"
	ProxyFile      string // proxy file
	BadgerDir      string // dir to store badger data
	BadgerValueDir string // dir to store badger value log in
	LogDir         string // dir to save log
}

// LoadConfig read config file, unmarshal data to config struct
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

// ProxyInfo proxy info
// set OnlyTCP true, if it only suports tcp,
// avoid of udp associate firstly Code,Addr,DNS must be not nil
type ProxyInfo struct {
	Code    string
	Addr    string
	User    string
	Pwd     string
	DNS     string
	OnlyTCP bool
}

// LoadProxies read proxy file, and unmarshal data
func LoadProxies(file string) (map[string]*ProxyInfo, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	proxies := make(map[string]*ProxyInfo)
	err = json.Unmarshal(buf, &proxies)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}
