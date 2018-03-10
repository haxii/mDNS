package main

import (
	"flag"

	"github.com/haxii/tdns"
)

var (
	port = flag.String("port", "8090", "port to listen on")
	dns  = flag.String("dns", "8.8.8.8:53", "dns name server")
	help = flag.Bool("h", false, "help usage")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	tDNS.SetDNSAddr(*dns)
	tDNS.Serve(":" + *port)
}
