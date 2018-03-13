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

	defaultRpcServer *RpcServer
	defaultLogger    *hlog.ZeroLogger
	defaultProxyMng  *proxy.ProxyManager
)

// Serve inits server and starts a rpc server
func Serve(configFile string) error {
	config, err := LoadConfig(configFile)
	if err != nil {
		return err
	}
	err = InitServer(config)
	if err != nil {
		return err
	}

	defaultRpcServer = &RpcServer{}
	err = defaultRpcServer.Start(config.ListenAddr)
	if err != nil {
		return err
	}
	return nil
}

// StopServer stops rpc server, closes db connect, resets proxy manager
func StopServer() {
	defaultRpcServer.Stop()
	geoip.CloseDB()
	badger.CloseDB()
	defaultProxyMng.Reset()
	defaultProxyMng = nil
	defaultRpcServer = nil
}

// InitServer inits db, logger, proxy manager
func InitServer(config *Config) error {
	// set dns server addr
	DNSServer = config.DNSServer

	// init logger
	var err error
	defaultLogger, err = hlog.MakeZeroLogger(false, config.LogDir, "tdns")
	if err != nil {
		return err
	}

	// init geoip db
	err = geoip.InitDB(config.IPDB)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}

	// init badger db
	err = badger.InitDB(config.BadgerDir, config.BadgerValueDir)
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}

	// config proxy manager
	defaultProxyMng = proxy.NewProxyManager()
	err = defaultProxyMng.LoadProxys()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	return nil
}

type RpcServer struct {
	rpcServer *gorpc.Server
}

// Start starts a rpc server
func (s *RpcServer) Start(addr string) error {
	s.rpcServer = NewRpcServer(addr)
	err := s.rpcServer.Start()
	if err != nil {
		defaultLogger.Error("server", err, "", "")
		return err
	}
	return nil
}

// Stop stops rpc server
func (s *RpcServer) Stop() {
	defaultRpcServer.rpcServer.Stop()
}
