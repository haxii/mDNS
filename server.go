package tdns

import (
	"encoding/json"
	"io/ioutil"
	"log"

	hlog "github.com/haxii/log"
	"github.com/haxii/tdns/db/badger"
	"github.com/haxii/tdns/db/geoip"
	"github.com/haxii/tdns/proxy"
	"github.com/valyala/gorpc"
)

var (
	DNSAddr = "8.8.8.8:53"

	defaultServer   *Server
	defaultLogger   *hlog.ZeroLogger
	defaultProxyMng *proxy.ProxyManager
	defaultConfig   *Config
)

type Server struct {
	//rpc server
	rpcServer *gorpc.Server
}

//Serve init server and listen on addr
func Serve(config string) {
	buf, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalln(err)
	}
	defaultConfig = &Config{}
	err = json.Unmarshal(buf, defaultConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defaultServer = &Server{}
	defaultServer.Init()
	defaultServer.Serve()
}

//Init
func (s *Server) Init() {
	//set dns server addr
	DNSAddr = defaultConfig.DNSServer

	var err error
	defaultLogger, err = hlog.MakeZeroLogger(false, defaultConfig.LogDir, "tdns")
	if err != nil {
		log.Fatalln(err)
	}

	//init geoip db
	err = geoip.InitDB(defaultConfig.IPDB)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}
	//init badger db
	err = badger.InitDB(defaultConfig.BadgerDir, defaultConfig.BadgerValueDir)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}

	//config proxy
	defaultProxyMng = proxy.NewProxyManager()
	err = defaultProxyMng.LoadProxys()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}
}

//Serve
func (s *Server) Serve() {
	s.rpcServer = NewRpcServer(defaultConfig.ListenAddr)
	err := s.rpcServer.Serve()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}
}
