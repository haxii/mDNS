package tdns

import (
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
func Serve(configFile string) error {
	config, err := LoadConfig(configFile)
	if err != nil {
		return err
	}
	defaultConfig = config
	err = Init()
	if err != nil {
		return err
	}

	defaultServer = &Server{}
	return defaultServer.Serve()
}

//StopServer stop server
func StopServer() {
	defaultServer.rpcServer.Stop()
	geoip.CloseDB()
	badger.CloseDB()
	defaultProxyMng.Reset()
	defaultProxyMng = nil
	defaultServer = nil
}

//Init
func Init() error {
	//set dns server addr
	DNSServer = defaultConfig.DNSServer

	//init logger
	var err error
	defaultLogger, err = hlog.MakeZeroLogger(false, defaultConfig.LogDir, "tdns")
	if err != nil {
		return err
	}

	//init geoip db
	err = geoip.InitDB(defaultConfig.IPDB)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	//init badger db
	err = badger.InitDB(defaultConfig.BadgerDir, defaultConfig.BadgerValueDir)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}

	//config proxy
	defaultProxyMng = proxy.NewProxyManager()
	err = defaultProxyMng.LoadProxys()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	return nil
}

//Serve
func (s *Server) Serve() error {
	s.rpcServer = NewRpcServer(defaultConfig.ListenAddr)
	err := s.rpcServer.Start()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	return nil
}
