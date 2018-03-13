package main

import (
	"flag"
	"log"

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
	err := tdns.Serve(*config)
	if err != nil {
		log.Fatalln(err)
	}
	<-stop
}
