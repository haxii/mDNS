package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/haxii/tdns"
)

var (
	rpc  = flag.String("rpc", "", "rpc server addr")
	help = flag.Bool("h", false, "help usage")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *rpc == "" {
		log.Fatalln("rpc is empty")
	}

	c := tdns.ConnectClient(*rpc)
	defer c.Close()
	proxys, err := c.ListProxyInfo()
	if err != nil {
		log.Fatalln(err)
	}
	for code, proxy := range proxys {
		fmt.Printf("%s: %#+v\n", code, proxy)
	}
}
