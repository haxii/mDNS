package proxy

import (
	"runtime"
	"testing"
)

func TestPoolGet(t *testing.T) {
	pool := socksClientPool{}
	socksClient := pool.Get()
	if socksClient == nil {
		t.Fail()
	}
}

func TestPoolPut(t *testing.T) {
	pool := socksClientPool{}
	socksClient := pool.Get()
	if socksClient == nil {
		t.Fail()
	}
	socksClient.Addr = "127.0.0.1:8000"
	pool.Put(socksClient)

	runtime.GC()
	socksClient = pool.Get()
	if socksClient.Addr != "" {
		t.Fatal("socksClient from pool is not empty")
	}
}
