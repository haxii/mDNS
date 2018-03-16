package tdns

import (
	"net"
	"strings"

	socks5 "github.com/nicdex/go-socks5"
)

//ProxyClient wraps socks client
type ProxyClient struct {
	dns                string
	client             *socks5.Client
	onlyTCP            bool //socks proxy only support tcp
	udpAssociateFailed int  //udp associate failed time
}

// NewProxyClient returns a proxy client
func NewProxyClient(addr, user, pwd, dns string, onlyTCP bool) *ProxyClient {
	client := &ProxyClient{
		client:  newSocksClient(addr, user, pwd),
		dns:     dns,
		onlyTCP: onlyTCP,
	}
	return client
}

// LookupIPAddrs send dns request, pars response
// return IPAddr slice and error if any
func (c *ProxyClient) LookupIPAddrs(host string) ([]net.IPAddr, error) {
	var conn net.Conn
	var err error
	if c.onlyTCP {
		conn, err = c.client.Dial("tcp", c.dns)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = c.client.Dial("udp", c.dns)
		if err == nil {
			c.udpAssociateFailed = 0
		} else {
			if !strings.Contains(err.Error(), "udp associate failed") {
				return nil, err
			}

			// socks proxy can't support udp associate, then use tcp to connect
			conn, err = c.client.Dial("tcp", c.dns)
			if err != nil {
				return nil, err
			}
			// if udp associate continuously fail three time, change onlyTCP true, use tcp next
			c.udpAssociateFailed++
			if c.udpAssociateFailed >= 3 {
				c.onlyTCP = true
			}
		}
	}
	defer conn.Close()

	ips, err := LookupIPOnConn(conn, host)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

// newSocksClient returns a socksClient
func newSocksClient(addr, user, pwd string) *socks5.Client {
	socksClient := &socks5.Client{
		Addr:     addr,
		Username: user,
		Password: pwd,
	}
	return socksClient
}
