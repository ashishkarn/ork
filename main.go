package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.Int("port", 9800, "udp/tcp port")
	flag.Parse()

	node := NewNode(*port)
	if err := node.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
