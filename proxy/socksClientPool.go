package proxy

import (
	"sync"

	socks5 "github.com/nicdex/go-socks5"
)

type socksClientPool struct {
	pool sync.Pool
}

//Get get a socksClient from pool
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

//Put put a socksClient to pool
func (p *socksClientPool) Put(c *socks5.Client) {
	c.Addr = ""
	c.Username = ""
	c.Password = ""
	c.Debug = false
	p.pool.Put(c)
}
