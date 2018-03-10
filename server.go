package tDNS

import (
	"github.com/haxii/tdns/proxy"
	"github.com/valyala/gorpc"
)

var (
	SocksAddr = "127.0.0.1:8000"
	DNSAddr   = "8.8.8.8:53"
)

var defaultServer *Server

type Server struct {
	//rpc server
	rpcServer *gorpc.Server
	//proxy manager
	proxyMng *proxy.ProxyManager
}

//init server and listen on addr
func Serve(addr string) {
	defaultServer = &Server{}
	defaultServer.Init()
	defaultServer.Serve(addr)
}

//Init
func (s *Server) Init() {
	//config proxy
	s.proxyMng = proxy.NewProxyManager()
	s.proxyMng.SetProxyClient(COUNTRY_CN, SocksAddr, "", "")
}

//Serve
func (s *Server) Serve(addr string) {
	s.rpcServer = NewRpcServer(addr)
	defaultServer.rpcServer.Serve()
}

//SetDNSAddr
func SetDNSAddr(dnsAddr string) {
	DNSAddr = dnsAddr
}
