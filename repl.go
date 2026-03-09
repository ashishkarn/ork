package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunREPL(node *Node) {
	fmt.Printf("ork running on :%d as %s\n", node.self.Port, node.self.ID)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		if len(args) == 0 {
			fmt.Print("> ")
			continue
		}

		switch args[0] {
		case "discover":
			peers := node.Discover()
			if len(peers) == 0 {
				fmt.Println("no peers found")
			} else {
				fmt.Printf("found %d peer(s):\n", len(peers))
				for _, p := range peers {
					fmt.Printf("  %s @ %s:%d\n", p.ID, p.Addr, p.Port)
				}
			}

		case "peers":
			peers := node.peers.All()
			if len(peers) == 0 {
				fmt.Println("no known peers")
			} else {
				for _, p := range peers {
					fmt.Printf("  %s @ %s:%d\n", p.ID, p.Addr, p.Port)
				}
			}

		case "quit", "exit":
			fmt.Println("bye")
			os.Exit(0)

		default:
			fmt.Printf("unknown command: %s\n", args[0])
		}

		fmt.Print("> ")
	}
}
