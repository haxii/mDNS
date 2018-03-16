package tdns

import (
	"testing"
)

func TestNewProxyClient(t *testing.T) {
	proxy := NewProxyClient(testAddr, testUser, testPwd, testDNS, true)
	if proxy.dns != testDNS {
		t.Fatalf("proxy dns(%s) != testDNS(%s)\n", proxy.dns, testDNS)
	}
	if proxy.onlyTCP != true {
		t.Fatal("proxy.onlyTCP not true")
	}
	if proxy.client == nil {
		t.Fatal("proxy.client nil")
	}
	if proxy.client.Addr != testAddr {
		t.Fatalf("proxy client addr(%s) != testAddr(%s)\n", proxy.client.Addr, testAddr)
	}
	if proxy.client.Username != testUser {
		t.Fatalf("proxy client Username(%s) != testUser(%s)\n", proxy.client.Username, testUser)
	}
	if proxy.client.Password != testPwd {
		t.Fatalf("proxy client Password(%s) != testPwd(%s)\n", proxy.client.Password, testPwd)
	}
}

func TestProxyLookupIPAddrs(t *testing.T) {
	proxyTCP := NewProxyClient(testAddr, testUser, testPwd, testDNS, true)
	ips, err := proxyTCP.LookupIPAddrs(testHost)
	if err != nil {
		t.Error(err)
	}
	if len(ips) == 0 {
		t.Fatal("no ip")
	}

	proxyUDP := NewProxyClient(testAddr, testUser, testPwd, testDNS, false)
	ips, err = proxyUDP.LookupIPAddrs(testHost)
	if err != nil {
		t.Error(err)
	}
	if len(ips) == 0 {
		t.Fatal("no ip")
	}
}
