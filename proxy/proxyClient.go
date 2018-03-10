package proxy

import (
	"net"

	"github.com/haxii/tdns/dns"
	socks5 "github.com/nicdex/go-socks5"
)

var (
	defaultSocksClientPool socksClientPool
)

type ProxyClient struct {
	client *socks5.Client
}

//Reset reset resource
func (c *ProxyClient) SetSocksClient(client *socks5.Client) {
	c.Reset()
	c.client = client
}

//newProxyClient new a socksClient
//return a ProxyClient
func newSocksClient(addr, user, pwd string) *socks5.Client {
	socksClient := defaultSocksClientPool.Get()
	socksClient.Addr = addr
	socksClient.Username = user
	socksClient.Password = pwd
	socksClient.Debug = false
	return socksClient
}

//ResoveDNS send dns request and parse response
//return IPAddr slice and error if any
func (c *ProxyClient) ResoveDNS(host, nameserver string) ([]net.IPAddr, error) {
	conn, err := c.client.Dial("udp", nameserver)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	ips, err := dns.LookupIPOnConn(conn, host)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

//Reset reset resource
func (c *ProxyClient) Reset() {
	if c.client != nil {
		defaultSocksClientPool.Put(c.client)
		c.client = nil
	}
}
