package tdns

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/haxii/tdns/db/badger"
)

var (
	testCountry = "CN"
	testAddr    = "127.0.0.1:8000"
	testUser    = ""
	testPwd     = ""
	testDNS     = "8.8.8.8:53"
	testIP      = "127.0.0.1"
	testHost    = "www.qq.com"

	testKey       = []byte("testKey")
	testBadgerDir = "./TestBadger"
	testIPs       = []net.IPAddr{net.IPAddr{
		IP: net.ParseIP(testIP),
	}}
)

func TestSetProxy(t *testing.T) {
	testTDNS := &TDNS{}
	testTDNS.SetProxy(testCountry, testAddr, testUser, testPwd,
		testDNS, true)
	val, ok := testTDNS.proxies.Load("CN")
	if !ok {
		t.Fatal("load proxy error")
	}
	proxy := val.(*ProxyClient)
	if proxy == nil {
		t.Fatal("proxy nil")
	}
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

func TestTDNSLookupIPAddrs(t *testing.T) {
	db, err := badger.OpenBadger(testBadgerDir, testBadgerDir)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	testTDNS := &TDNS{
		BadgerDB: db,
	}
	testTDNS.SetProxy(testCountry, testAddr, testUser, testPwd,
		testDNS, true)

	ips, err := testTDNS.LookupIPAddrs(testCountry, testHost)
	if err != nil {
		t.Error(err)
	}
	if len(ips) == 0 {
		t.Fatal("no ip")
	}
	os.RemoveAll(testBadgerDir)
}

func TestSaveIPsToCache(t *testing.T) {
	db, err := badger.OpenBadger(testBadgerDir, testBadgerDir)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	testTDNS := &TDNS{
		BadgerDB: db,
		CacheTTL: time.Second,
	}

	testTDNS.saveIPsToCache(testKey, testIPs)
	bs, err := testTDNS.BadgerDB.Get(testKey)
	if err != nil {
		t.Error(err)
	}
	if bs != nil {
		ips := make([]net.IPAddr, 1)
		err = json.Unmarshal(bs, &ips)
		if err != nil {
			t.Error(err)
		}
		if len(ips) != 1 {
			t.Fatal("ips length != 1")
		}
		if ips[0].IP.String() != testIP {
			t.Fatalf("ip(%s) != %s", ips[0].IP.String(), testIP)
		}
	}

	time.Sleep(time.Second)
	_, err = testTDNS.BadgerDB.Get(testKey)
	if err != badger.ErrKeyNotFound {
		t.Fatalf("after ttl, but find key")
	}

	os.RemoveAll(testBadgerDir)
}
