package main

import (
	"errors"
	"net"

	"github.com/valyala/gorpc"
)

// NewRpcServer makes a tcp rpc server
func NewRpcServer(addr string) *gorpc.Server {
	d := gorpc.NewDispatcher()
	d.AddFunc("LookupIPAddr", LookupIPAddrsForServer)
	return gorpc.NewTCPServer(addr, d.NewHandlerFunc())
}

type LookupIPRequest struct {
	Code string
	Host string
}

// LookupIPAddrs returns ip slice and error if any
func LookupIPAddrsForServer(req *LookupIPRequest) ([]net.IPAddr, error) {
	if len(req.Code) == 0 || len(req.Host) == 0 {
		return nil, errors.New("ip or host is empty")
	}
	return defaultTDNS.LookupIPAddrs(req.Code, req.Host)
}

// LookupIPAddrForClient
//
// new a tcp client and call rpc
// returns IPAddr slice and err if any
func LookupIPAddrForClient(rpc, code, domain string) ([]net.IPAddr, error) {
	rpcClient := gorpc.NewTCPClient(rpc)
	rpcClient.Start()
	defer rpcClient.Stop()

	d := gorpc.NewDispatcher()
	d.AddFunc("LookupIPAddr", LookupIPAddrsForServer)
	dc := d.NewFuncClient(rpcClient)
	req := &LookupIPRequest{
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
