package proxy

import (
	"net"
	"strings"

	"github.com/haxii/tdns/dns"
	socks5 "github.com/nicdex/go-socks5"
)

var (
	defaultSocksClientPool socksClientPool
)

type ProxyClient struct {
	client             *socks5.Client
	onlyTCP            bool //socks proxy only support tcp
	udpAssociateFailed int  //udp associate failed time
}

// Reset resets resource
func (c *ProxyClient) SetSocksClient(client *socks5.Client) {
	c.Reset()
	c.client = client
}

// SetOnlyTCP sets onlyTCP
func (c *ProxyClient) SetOnlyTCP(onlyTCP bool) {
	c.onlyTCP = onlyTCP
}

// NewSocksClient returns a socksClient
func NewSocksClient(addr, user, pwd string) *socks5.Client {
	socksClient := defaultSocksClientPool.Get()
	socksClient.Addr = addr
	socksClient.Username = user
	socksClient.Password = pwd
	socksClient.Debug = false
	return socksClient
}

// ResoveDNS sends dns request, parses response
// ResoveDNS returns IPAddr slice and error if any
func (c *ProxyClient) ResoveDNS(host, dnsServer string) ([]net.IPAddr, error) {
	var conn net.Conn
	var err error
	if c.onlyTCP {
		conn, err = c.client.Dial("tcp", dnsServer)
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = c.client.Dial("udp", dnsServer)
		if err != nil {
			if !strings.Contains(err.Error(), "udp associate failed") {
				return nil, err
			}

			// socks proxy can't support udp associate, then use tcp to connect
			conn, err = c.client.Dial("tcp", dnsServer)
			if err != nil {
				return nil, err
			}
			// if udp associate continuously fail on three time, change onlyTCP true on memory
			c.udpAssociateFailed++
			if c.udpAssociateFailed >= 3 {
				c.onlyTCP = true
			}
		} else {
			c.udpAssociateFailed = 0
		}
	}
	defer conn.Close()

	ips, err := dns.LookupIPOnConn(conn, host)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

// Reset resets resource
func (c *ProxyClient) Reset() {
	if c.client != nil {
		defaultSocksClientPool.Put(c.client)
		c.client = nil
	}
}
