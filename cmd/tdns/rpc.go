package main

import (
	"errors"
	"net"

	"github.com/valyala/gorpc"
)

// LookupIPService LookupIP service
type LookupIPService struct {
}

type lookupIPRequest struct {
	Code string
	Host string
}

//LookupIPAddr returns ip slice and error if any
func (s *LookupIPService) LookupIPAddr(req *lookupIPRequest) ([]net.IPAddr, error) {
	if len(req.Code) == 0 || len(req.Host) == 0 {
		return nil, errors.New("code or domain is empty")
	}
	return defaultTDNS.LookupIPAddrs(req.Code, req.Host, defaultCode)
}

// NewRPCServer makes a tcp rpc server
func NewRPCServer(addr string) *gorpc.Server {
	d := gorpc.NewDispatcher()
	service := &LookupIPService{}
	d.AddService("LookupIPAddr", service)
	return gorpc.NewTCPServer(addr, d.NewHandlerFunc())
}

// LookupIPAddrForClient new a tcp client and call rpc
//
// returns IPAddr slice and err if any
func LookupIPAddrForClient(rpc, code, domain string) ([]net.IPAddr, error) {
	rpcClient := gorpc.NewTCPClient(rpc)
	rpcClient.Start()
	defer rpcClient.Stop()

	d := gorpc.NewDispatcher()
	service := &LookupIPService{}
	d.AddService("LookupIPAddr", service)

	dc := d.NewServiceClient("LookupIPAddr", rpcClient)
	req := &lookupIPRequest{
		Code: code,
		Host: domain,
	}
	resp, err := dc.Call("LookupIPAddr", req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("response nil")
	}
	ipAddrs := resp.([]net.IPAddr)
	return ipAddrs, nil
}
