package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/haxii/daemon"
	"github.com/haxii/log"
	"github.com/haxii/tdns"
	"github.com/haxii/tdns/db/badger"
)

var (
	defaultLogger *log.ZeroLogger
	defaultTDNS   *tdns.TDNS
)

var (
	config  = flag.String("config", "", "server config")
	rpc     = flag.String("rpc", "", "rpc addr")
	resolve = flag.String("resolve", "", "host to resolve")
	country = flag.String("country", "", "country code")
	_       = flag.String("s", daemon.UsageDefaultName, daemon.UsageMessage)
)

func main() {
	flag.Parse()

	// resolve with country code
	if len(*rpc) > 0 && len(*resolve) > 0 && len(*country) > 0 {
		ips, err := LookupIPAddrForClient(*rpc, *country, *resolve)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s ---> %v\n", *resolve, ips)
		return
	}

	daemon.Make("-s", "tdns", "tdns daemon service").Run(func() {
		if len(*config) > 0 {
			serve(*config)
		}
	})
}

func serve(configFile string) {
	config, err := LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	// init logger
	defaultLogger, err = log.MakeZeroLogger(false, config.LogDir, "tdns")
	if err != nil {
		panic(err)
	}

	// init badger db
	badgerDB, err := badger.OpenBadger(config.BadgerDir, config.BadgerValueDir)
	if err != nil {
		panic(err)
	}

	defaultTDNS = &tdns.TDNS{
		Logger:   defaultLogger,
		BadgerDB: badgerDB,
		CacheTTL: time.Hour * 72,
	}

	// load proxy info
	proxies, err := LoadProxies(config.ProxyFile)
	if err != nil {
		panic(err)
	}

	for _, proxy := range proxies {
		defaultTDNS.SetProxy(proxy.Code, proxy.Addr, proxy.User,
			proxy.Pwd, proxy.DNS, proxy.OnlyTCP)
	}

	// start rcp server
	rpcServer := NewRpcServer(config.RpcAddr)
	fmt.Printf("listen rpc on: %s\n", config.RpcAddr)
	err = rpcServer.Serve()
	if err != nil {
		panic(err)
	}
}
