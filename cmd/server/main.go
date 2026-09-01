package main

import (
	"flag"

	"github.com/Meh-Mehul/sf/server"
)

func main() {
	hubFlag := flag.String("hub", "localhost:9999", "central hub address to connect to")
	flag.Parse()

	opts := &server.ServerOpts{
		Address:    *hubFlag,
		MaxSession: 1,
	}

	server.StartServer(opts)
}
