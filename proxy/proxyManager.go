package proxy

import (
	"sync"
)

type ProxyManager struct {
	//CountryCode:ProxyClient
	proxys map[string]*ProxyClient

	sync.RWMutex
}

//NewProxyManager create a new ProxyManager
func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		proxys: make(map[string]*ProxyClient),
	}
}

//GetProxyClient return ProxyClient
func (m *ProxyManager) GetProxyClient(code string) *ProxyClient {
	m.RLock()
	defer m.RUnlock()
	return m.proxys[code]
}

//SetProxyClient set new socks client for a country code
func (m *ProxyManager) SetProxyClient(code, addr, user, pwd string) {
	m.Lock()
	defer m.Unlock()
	socks := newSocksClient(addr, user, pwd)
	if socks != nil {
		if m.proxys[code] == nil {
			m.proxys[code] = &ProxyClient{}
		}
		m.proxys[code].SetSocksClient(socks)
	}
}
