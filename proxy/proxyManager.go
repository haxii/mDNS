package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/haxii/tdns/db/badger"
)

var (
	countryCodeKey     = []byte("country_codes")
	countrySeparator   = []byte(":")
	proxyInfoKeyPrefix = []byte("proxy:")
)

type ProxyInfo struct {
	Addr    string // socks addr
	User    string // socks user
	Pwd     string // socks password
	OnlyTCP bool   // only support tcp, default false
}

type ProxyManager struct {
	countryCodes [][]byte
	//CountryCode:ProxyClient
	proxys map[string]*ProxyClient

	sync.RWMutex
}

// NewProxyManager creates a new ProxyManager
func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		proxys:       make(map[string]*ProxyClient),
		countryCodes: make([][]byte, 0),
	}
}

// LoadProxys loads proxys from db
func (m *ProxyManager) LoadProxys() error {
	val, err := badger.Get(countryCodeKey)
	if err != nil {
		if err.Error() != "Key not found" {
			return err
		} else {
			return nil
		}
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
		socks := NewSocksClient(info.Addr, info.User, info.Pwd)
		if socks != nil {
			codeStr := string(code)
			if m.proxys[codeStr] == nil {
				m.proxys[codeStr] = &ProxyClient{}
			}
			m.proxys[codeStr].SetSocksClient(socks)
			m.proxys[codeStr].SetOnlyTCP(info.OnlyTCP)
		}
	}

	return nil
}

// GetProxyClient returns ProxyClient
func (m *ProxyManager) GetProxyClient(code string) *ProxyClient {
	m.RLock()
	defer m.RUnlock()
	return m.proxys[code]
}

// SetProxyClient saves proxy info to db and sets socks client
func (m *ProxyManager) SetProxy(code, addr, user, pwd string, onlyTCP bool) error {
	//save to db
	info := &ProxyInfo{
		Addr:    addr,
		User:    user,
		Pwd:     pwd,
		OnlyTCP: onlyTCP,
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
	socks := NewSocksClient(addr, user, pwd)
	if socks != nil {
		m.Lock()
		if m.proxys[code] == nil {
			m.proxys[code] = &ProxyClient{}
		}
		m.proxys[code].SetSocksClient(socks)
		m.proxys[code].SetOnlyTCP(onlyTCP)
		m.Unlock()
	}
	return nil
}

// SetProxyOnlyTCP sets proxy onlyTCP info
func (m *ProxyManager) SetProxyOnlyTCP(code string, onlyTCP bool) error {
	bs, err := badger.Get(getProxyKey([]byte(code)))
	if err != nil {
		return err
	}

	info := &ProxyInfo{}
	err = json.Unmarshal(bs, info)
	if err != nil {
		return err
	}

	info.OnlyTCP = onlyTCP
	err = badger.Set(getProxyKey([]byte(code)), bs)
	if err != nil {
		return err
	}

	bs, err = json.Marshal(info)
	if err != nil {
		return err
	}
	err = badger.Set(getProxyKey([]byte(code)), bs)
	if err != nil {
		return err
	}

	if m.proxys[code] == nil {
		return errors.New("not found proxy in proxy manager")
	}
	m.Lock()
	m.proxys[code].SetOnlyTCP(onlyTCP)
	m.Unlock()

	return nil
}

// GetProxys returns proxy map
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

// Reset resets proxy client, countryCodes
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
