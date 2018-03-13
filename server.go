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
	DNSServer = "8.8.8.8:53"

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
func Serve(config string) error {
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

	return defaultServer.Serve()
}

func StopServer() {
	defaultServer.rpcServer.Stop()
	geoip.CloseDB()
	badger.CloseDB()
	defaultProxyMng.Reset()
	defaultProxyMng = nil
	defaultServer = nil
}

//Init
func (s *Server) Init() {
	//set dns server addr
	DNSServer = defaultConfig.DNSServer

	var err error
	defaultLogger, err = hlog.MakeZeroLogger(false, defaultConfig.LogDir, "tdns")
	if err != nil {
		log.Fatalln(err)
	}

	//init geoip db
	err = geoip.InitDB(defaultConfig.IPDB)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		log.Fatalln(err)
	}
	//init badger db
	err = badger.InitDB(defaultConfig.BadgerDir, defaultConfig.BadgerValueDir)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		log.Fatalln(err)
	}

	//config proxy
	defaultProxyMng = proxy.NewProxyManager()
	err = defaultProxyMng.LoadProxys()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}
}

//Serve
func (s *Server) Serve() error {
	s.rpcServer = NewRpcServer(defaultConfig.ListenAddr)
	err := s.rpcServer.Start()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
	}
	return err
}
