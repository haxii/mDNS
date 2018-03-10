package main

import (
	"fmt"
	"log"

	"github.com/haxii/tdns"
)

func main() {
	c := tDNS.ConnectClient("127.0.0.1:8090")
	defer c.Close()

	var name = "www.google.com"
	res, err := c.LookupIPAddr("61.135.169.125", name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s-->%+v\n", name, res)
}
