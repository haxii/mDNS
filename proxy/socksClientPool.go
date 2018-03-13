package proxy

import (
	"sync"

	socks5 "github.com/nicdex/go-socks5"
)

type socksClientPool struct {
	pool sync.Pool
}

// Get gets a socksClient from pool
func (p *socksClientPool) Get() *socks5.Client {
	v := p.pool.Get()
	var client *socks5.Client
	if v != nil {
		client = v.(*socks5.Client)
	} else {
		client = &socks5.Client{}
	}
	return client
}

//Put puts a socksClient back to pool
func (p *socksClientPool) Put(c *socks5.Client) {
	p.pool.Put(c)
}
