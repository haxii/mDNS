package main

import (
	"flag"
	"log"

	"github.com/haxii/tdns"
)

var (
	rpc     = flag.String("rpc", "", "rpc server addr")
	code    = flag.String("code", "", "country code")
	addr    = flag.String("addr", "", "socks addr, host with port")
	user    = flag.String("user", "", "username")
	pwd     = flag.String("pwd", "", "password")
	onlyTCP = flag.Bool("onlyTCP", false, "only support tcp")
	help    = flag.Bool("h", false, "help usage")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *rpc == "" || *code == "" {
		log.Fatalln("rpc or code or addr is empty")
	}

	c := tdns.ConnectClient(*rpc)
	defer c.Close()

	err := c.SetProxyInfo(*code, *addr, *user, *pwd, *onlyTCP)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("ok")
}
