package proxy

import (
	"os"
	"testing"

	"github.com/haxii/tdns/db/badger"
)

var (
	badgerDir = "./badger"
	mng       *ProxyManager
)

func TestProxyManager(t *testing.T) {
	badger.InitDB(badgerDir, badgerDir)

	t.Run("testNewProxyManager", testNewProxyManager)
	t.Run("testSetProxy", testSetProxy)
	t.Run("testGetProxyClient", testGetProxyClient)
	t.Run("testLoadProxys", testLoadProxys)
	t.Run("testGetProxys", testGetProxys)

	os.RemoveAll(badgerDir)
}

func testNewProxyManager(t *testing.T) {
	mng = NewProxyManager()
	if mng == nil || mng.countryCodes == nil || mng.proxys == nil {
		t.Fail()
	}
}

func testSetProxy(t *testing.T) {
	err := mng.SetProxy(testCountryCode, testProxyAddr, "user", "pwd", false)
	if err != nil {
		t.Error(err)
	}
}

func testGetProxyClient(t *testing.T) {
	client := mng.GetProxyClient(testCountryCode)
	if client == nil {
		t.Fail()
	}
}

func testLoadProxys(t *testing.T) {
	err := mng.LoadProxys()
	if err != nil {
		t.Error(err)
	}
}

func testGetProxys(t *testing.T) {
	proxys, err := mng.GetProxys()
	if err != nil {
		t.Error(err)
	}
	if len(proxys) != 1 {
		t.Fail()
	}
	if proxys[testCountryCode] == nil {
		t.Fail()
	}
}
