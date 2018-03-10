package main

import (
	"flag"
	"log"

	"github.com/haxii/tdns"
)

var (
	rpc  = flag.String("rpc", "", "rpc server addr")
	code = flag.String("code", "", "country code")
	addr = flag.String("addr", "", "socks addr, host with port")
	user = flag.String("user", "", "username")
	pwd  = flag.String("pwd", "", "password")
	help = flag.Bool("h", false, "help usage")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *rpc == "" || *code == "" || *addr == "" {
		log.Fatalln("rpc or code or addr is empty")
	}

	c := tDNS.ConnectClient(*rpc)
	defer c.Close()

	err := c.SetProxyInfo(*code, *addr, *user, *pwd)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("ok")
}
