package main

import (
	"flag"

	"github.com/haxii/tdns"
)

var (
	config = flag.String("config", "./file/config.json", "server config")
	help   = flag.Bool("h", false, "help usage")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	stop := make(chan struct{})
	tdns.Serve(*config)
	<-stop
}
