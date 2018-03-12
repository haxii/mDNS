package tdns

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/haxii/tdns/db/badger"
	"github.com/haxii/tdns/db/geoip"
	"github.com/haxii/tdns/proxy"
	"github.com/valyala/gorpc"
)

var (
	DNSAddr = "8.8.8.8:53"
)

var defaultServer *Server

type Server struct {
	//rpc server
	rpcServer *gorpc.Server
	//proxy manager
	proxyMng *proxy.ProxyManager
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

	//init geoip db
	err := geoip.InitDB(defaultConfig.IPDB)
	if err != nil {
		log.Fatalln(err)
	}
	//init badger db
	err = badger.InitDB(defaultConfig.BadgerDir, defaultConfig.BadgerValueDir)
	if err != nil {
		log.Fatalln(err)
	}

	//config proxy
	s.proxyMng = proxy.NewProxyManager()
	s.proxyMng.LoadProxys()
}

//Serve
func (s *Server) Serve() {
	s.rpcServer = NewRpcServer(defaultConfig.ListenAddr)
	defaultServer.rpcServer.Serve()
}
