package proxy

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/haxii/tdns/db/badger"
)

var (
	countryCodeKey     = []byte("country_codes")
	countrySeparator   = []byte(":")
	proxyInfoKeyPrefix = []byte("proxy:")
)

type ProxyInfo struct {
	Addr string
	User string
	Pwd  string
}

type ProxyManager struct {
	countryCodes [][]byte
	//CountryCode:ProxyClient
	proxys map[string]*ProxyClient

	sync.RWMutex
}

//NewProxyManager create a new ProxyManager
func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		proxys:       make(map[string]*ProxyClient),
		countryCodes: make([][]byte, 0),
	}
}

//LoadProxys load proxys from db
func (m *ProxyManager) LoadProxys() error {
	val, err := badger.Get(countryCodeKey)
	if err != nil {
		return err
	}

	codes := bytes.Split(val, countrySeparator)
	m.countryCodes = codes
	for _, code := range codes {
		val, err := badger.Get(getProxyKey(code))
		if err != nil {
			return err
		}
		info := &ProxyInfo{}
		err = json.Unmarshal(val, info)
		if err != nil {
			return err
		}
		socks := newSocksClient(info.Addr, info.User, info.Pwd)
		if socks != nil {
			codeStr := string(code)
			if m.proxys[codeStr] == nil {
				m.proxys[codeStr] = &ProxyClient{}
			}
			m.proxys[codeStr].SetSocksClient(socks)
		}
	}

	return nil
}

//GetProxyClient return ProxyClient
func (m *ProxyManager) GetProxyClient(code string) *ProxyClient {
	m.RLock()
	defer m.RUnlock()
	return m.proxys[code]
}

//SetProxyClient save proxy info to db and new socks client
func (m *ProxyManager) SetProxy(code, addr, user, pwd string) error {
	//save to db
	info := &ProxyInfo{
		Addr: addr,
		User: user,
		Pwd:  pwd,
	}
	bs, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = badger.Set(getProxyKey([]byte(code)), bs)
	if err != nil {
		return err
	}

	if m.proxys[code] == nil {
		m.countryCodes = append(m.countryCodes, []byte(code))
		err = badger.Set(countryCodeKey, bytes.Join(m.countryCodes, countrySeparator))
		if err != nil {
			return err
		}
	}

	//set socks client
	socks := newSocksClient(addr, user, pwd)
	if socks != nil {
		m.Lock()
		if m.proxys[code] == nil {
			m.proxys[code] = &ProxyClient{}
		}
		m.proxys[code].SetSocksClient(socks)
		m.Unlock()
	}
	return nil
}

//GetProxys return proxy map
func (m *ProxyManager) GetProxys() (map[string]*ProxyInfo, error) {
	val, err := badger.Get(countryCodeKey)
	if err != nil {
		return nil, err
	}

	codes := bytes.Split(val, countrySeparator)
	proxyInfos := make(map[string]*ProxyInfo)
	for _, code := range codes {
		val, err := badger.Get(getProxyKey(code))
		if err != nil {
			return nil, err
		}
		info := &ProxyInfo{}
		err = json.Unmarshal(val, info)
		if err != nil {
			return nil, err
		}
		proxyInfos[string(code)] = info
	}

	return proxyInfos, nil
}

//Reset reset proxy client, countryCodes
func (m *ProxyManager) Reset() {
	m.Lock()
	defer m.Unlock()

	for _, client := range m.proxys {
		client.Reset()
	}
	m.proxys = make(map[string]*ProxyClient)
	m.countryCodes = m.countryCodes[:0]
}

func getProxyKey(code []byte) []byte {
	return append(proxyInfoKeyPrefix, []byte(code)...)
}
