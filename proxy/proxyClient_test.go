package proxy

import (
	"testing"

	socks5 "github.com/nicdex/go-socks5"
)

var (
	testCountryCode = "CN"
	testProxyAddr   = "127.0.0.1:8000"
	testHost        = "www.qq.com"
)

var (
	socksClient *socks5.Client
	proxyClient *ProxyClient
)

func TestProxyClient(t *testing.T) {
	t.Run("testNewSocksClient", testNewSocksClient)
	t.Run("testSetSocksClient", testSetSocksClient)
	//t.Run("testResoveDNS", testResoveDNS)
	t.Run("testProxyClietReset", testProxyClietReset)
}

func testNewSocksClient(t *testing.T) {
	addr := testProxyAddr
	user := ""
	pwd := ""
	socksClient = NewSocksClient(addr, user, pwd)
	if socksClient.Addr != addr {
		t.Fail()
	}
	if socksClient.Username != user {
		t.Fail()
	}
	if socksClient.Password != pwd {
		t.Fail()
	}
}

func testSetSocksClient(t *testing.T) {
	proxyClient = &ProxyClient{}
	proxyClient.SetSocksClient(socksClient)
	if proxyClient.client == nil {
		t.Fail()
	}
}

/*func testResoveDNS(t *testing.T) {
	ips, err := proxyClient.ResoveDNS(testHost, "8.8.8.8:53")
	if err != nil {
		t.Error(err)
	}
	if ips == nil {
		t.Fail()
	}
}*/

func testProxyClietReset(t *testing.T) {
	proxyClient.Reset()
	if proxyClient.client != nil {
		t.Fail()
	}
}
