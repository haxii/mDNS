package tdns

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/haxii/log"
	"github.com/haxii/tdns/db/badger"
)

type TDNS struct {
	Logger   log.Logger
	BadgerDB *badger.BadgerDB
	CacheTTL time.Duration

	proxies sync.Map
}

// SetProxy
//
// new a proxy client and store it on proxies
func (tdns *TDNS) SetProxy(code, addr, user, pwd, dns string, onlyTCP bool) {
	proxy := NewProxyClient(addr, user, pwd, dns, onlyTCP)
	tdns.proxies.Store(code, proxy)
}

// LookupIPAddrs
//
// read from cache, if no cache, then resolve it and save in cache async
func (tdns *TDNS) LookupIPAddrs(code, host string) ([]net.IPAddr, error) {
	if len(code) == 0 || len(host) == 0 {
		return nil, errors.New("code or host is empty")
	}

	var cacheKey []byte
	if tdns.BadgerDB != nil {
		//get ips from cache
		cacheKey = []byte(fmt.Sprintf("host_ips:%s:%s", host, code))
		bs, err := tdns.BadgerDB.Get(cacheKey)
		if err != nil && err != badger.ErrKeyNotFound {
			return nil, err
		}
		if bs != nil {
			ips := make([]net.IPAddr, 1)
			err = json.Unmarshal(bs, &ips)
			if err != nil {
				return nil, err
			}
			return ips, nil
		}
	}

	//resolve on proxy
	value, ok := tdns.proxies.Load(code)
	if !ok {
		return nil, errors.New("not found proxy")
	}
	proxy := value.(*ProxyClient)
	ips, err := proxy.LookupIPAddrs(host)
	if err != nil {
		return nil, err
	}

	// async save data to cache
	go tdns.saveIPsToCache(cacheKey, ips)

	return ips, err
}

// save data to cache if BadgerDB not nil
func (tdns *TDNS) saveIPsToCache(key []byte, ips []net.IPAddr) {
	defer func() {
		if r := recover(); r != nil {
			tdns.Logger.Error("TDNS", nil, "%v\n%s", r, debug.Stack())
		}
	}()

	if tdns.BadgerDB == nil {
		return
	}

	bs, err := json.Marshal(ips)
	if err != nil {
		tdns.Logger.Error("TDNS", err, "ips: %#+v", ips)
		return
	}
	if tdns.CacheTTL > 0 {
		err = tdns.BadgerDB.SetWithTTL(key, bs, tdns.CacheTTL)
		if err != nil {
			tdns.Logger.Error("TDNS", err, "", "")
		}
	} else {
		err = tdns.BadgerDB.Set(key, bs)
		if err != nil {
			tdns.Logger.Error("TDNS", err, "", "")
		}
	}
}
