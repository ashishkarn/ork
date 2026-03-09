package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ork <daemon|discover|stop>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "daemon":
		fs := flag.NewFlagSet("daemon", flag.ExitOnError)
		port := fs.Int("port", 9800, "udp/tcp port")
		fs.Parse(os.Args[2:])
		d := NewDaemon(*port)
		if err := d.Start(); err != nil {
			log.Fatal(err)
		}

	case "discover":
		resp, err := daemonCmd(CmdDiscover)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.Peers) == 0 {
			fmt.Println("no peers found")
			return
		}
		fmt.Printf("found %d peer(s):\n", len(resp.Peers))
		for _, p := range resp.Peers {
			fmt.Printf("  %s @ %s:%d\n", p.ID, p.Addr, p.Port)
		}

	case "stop":
		if _, err := daemonCmd(CmdStop); err != nil {
			log.Fatal(err)
		}
		fmt.Println("daemon stopped")

	default:
		fmt.Println("usage: ork <daemon|discover|stop>")
		os.Exit(1)
	}
}
