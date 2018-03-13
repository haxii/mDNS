package tdns

import (
	"errors"
	"net"

	"github.com/haxii/tdns/proxy"
	"github.com/valyala/gorpc"
)

type Client struct {
	c  *gorpc.Client
	dc *gorpc.DispatcherClient
}

// ConnectClient connects rpc server
// ConnectClient returns a rpc client
func ConnectClient(addr string) *Client {
	c := gorpc.NewTCPClient(addr)
	c.Start()

	d := gorpc.NewDispatcher()
	d.AddFunc("LookupIPAddr", LookupIPAddrs)
	d.AddFunc("SetProxyInfo", SetProxyInfo)
	d.AddFunc("ListProxyInfo", ListProxyInfo)
	dc := d.NewFuncClient(c)
	client := &Client{
		c:  c,
		dc: dc,
	}
	return client
}

// Close close client
func (c *Client) Close() {
	c.c.Stop()
}

// LookupIPAddr calls "LookupIPAddr" rpc
// LookupIPAddr returns IPAddr slice and err
func (c *Client) LookupIPAddr(ip, domain string) ([]net.IPAddr, error) {
	req := &LookupIPRequest{
		IP:   ip,
		Host: domain,
	}
	resp, err := c.dc.Call("LookupIPAddr", req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("response nil")
	}
	ipAddrs := resp.([]net.IPAddr)
	return ipAddrs, nil
}

// SetProxyInfo calls "SetProxyInfo" rpc
func (c *Client) SetProxyInfo(code, addr, user, pwd string, onlyTCP bool) error {
	req := &SetProxyRequest{
		Code:    code,
		Addr:    addr,
		User:    user,
		Pwd:     pwd,
		OnlyTCP: onlyTCP,
	}
	_, err := c.dc.Call("SetProxyInfo", req)
	if err != nil {
		return err
	}
	return nil
}

// ListProxyInfo calls "ListProxyInfo" rpc
// ListProxyInfo returns proxy map, error if any
func (c *Client) ListProxyInfo() (map[string]*proxy.ProxyInfo, error) {
	resp, err := c.dc.Call("ListProxyInfo", nil)
	if resp != nil {
		return resp.(map[string]*proxy.ProxyInfo), err
	}
	return nil, err
}
