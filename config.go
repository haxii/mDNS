package tdns

var defaultConfig *Config

type Config struct {
	ListenAddr     string // listen addr, eg. ":8080", "127.0.0.1:8080"
	DNSServer      string // dns server addr, eg. "8.8.8.8:53"
	IPDB           string // ip db file
	BadgerDir      string // dir to store badger data
	BadgerValueDir string // dir to stroe badger value log in
}
