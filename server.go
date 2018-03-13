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
)

//Serve init server and listen on addr
func Serve(configFile string) error {
	config, err := LoadConfig(configFile)
	if err != nil {
		return err
	}
	err = InitServer(config)
	if err != nil {
		return err
	}

	defaultServer = &Server{}
	return defaultServer.Start(config.ListenAddr)
}

//StopServer stop server
func StopServer() {
	defaultServer.Stop()
	geoip.CloseDB()
	badger.CloseDB()
	defaultProxyMng.Reset()
	defaultProxyMng = nil
	defaultServer = nil
}

//InitServer
func InitServer(config *Config) error {
	//set dns server addr
	DNSServer = config.DNSServer

	//init logger
	var err error
	defaultLogger, err = hlog.MakeZeroLogger(false, config.LogDir, "tdns")
	if err != nil {
		return err
	}

	//init geoip db
	err = geoip.InitDB(config.IPDB)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	//init badger db
	err = badger.InitDB(config.BadgerDir, config.BadgerValueDir)
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

type Server struct {
	//rpc server
	rpcServer *gorpc.Server
}

//Start
func (s *Server) Start(addr string) error {
	s.rpcServer = NewRpcServer(addr)
	err := s.rpcServer.Start()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	return nil
}

//Stop
func (s *Server) Stop() {
	defaultServer.rpcServer.Stop()
}
