package main

import (
	"flag"
	"log"

	"github.com/Meh-Mehul/sf/central"
)

func main() {
	addrFlag := flag.String("addr", ":9999", "listen address for central hub")
	flag.Parse()

	opts := &central.CentralOpts{
		Address: *addrFlag,
	}

	if err := central.StartHub(opts); err != nil {
		log.Fatalf("He died: %v\n", err)
	}
}
