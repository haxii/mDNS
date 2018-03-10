package tDNS

import (
	"errors"
	"net"

	"github.com/haxii/tdns/model/rpc"
	"github.com/valyala/gorpc"
)

func NewRpcServer(addr string) *gorpc.Server {
	d := gorpc.NewDispatcher()
	d.AddFunc("LookupIPAddr", LookupIPAddrs)
	d.AddFunc("SetProxyInfo", SetProxyInfo)
	return gorpc.NewTCPServer(addr, d.NewHandlerFunc())
}

//LookupIPAddrs rpc function
func LookupIPAddrs(req *rpc.LookupIPRequest) ([]net.IPAddr, error) {
	if len(req.IP) == 0 || len(req.Host) == 0 {
		return nil, errors.New("ip or host is empty")
	}
	countryCode := getCoutryCodeByIP(req.IP)
	client := defaultServer.proxyMng.GetProxyClient(countryCode)
	return client.ResoveDNS(req.Host, DNSAddr)
}

//SetProxyInfo rpc function
func SetProxyInfo(req *rpc.SetProxyRequest) error {
	if len(req.Code) == 0 || len(req.Addr) == 0 {
		return errors.New("code or addr is empty")
	}
	defaultServer.proxyMng.SetProxyClient(req.Code, req.Addr, req.User, req.Pwd)
	return nil
}

//return country code
func getCoutryCodeByIP(ip string) string {
	return COUNTRY_CN
}
