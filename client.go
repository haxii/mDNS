package tDNS

import (
	"errors"
	"net"

	"github.com/haxii/tdns/model/rpc"
	"github.com/valyala/gorpc"
)

type Client struct {
	c  *gorpc.Client
	dc *gorpc.DispatcherClient
}

//connect rpc server
//return rpc client
func ConnectClient(addr string) *Client {
	c := gorpc.NewTCPClient(addr)
	c.Start()

	d := gorpc.NewDispatcher()
	d.AddFunc("LookupIPAddr", LookupIPAddrs)
	d.AddFunc("SetProxyInfo", SetProxyInfo)
	dc := d.NewFuncClient(c)
	client := &Client{
		c:  c,
		dc: dc,
	}
	return client
}

//Close close client
func (c *Client) Close() {
	c.c.Stop()
}

//LookupIPAddr call "LookupIPAddr" rpc
//return IPAddr slice and err
func (c *Client) LookupIPAddr(ip, domain string) ([]net.IPAddr, error) {
	req := &rpc.LookupIPRequest{
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

//SetProxyInfo call "SetProxyInfo" rpc
func (c *Client) SetProxyInfo(code, addr, user, pwd string) error {
	req := &rpc.SetProxyRequest{
		Code: code,
		Addr: addr,
		User: user,
		Pwd:  pwd,
	}
	_, err := c.dc.Call("SetProxyInfo", req)
	if err != nil {
		return err
	}
	return nil
}
