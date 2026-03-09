package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	sockPath = "/tmp/ork.sock"
	pidPath  = "/tmp/ork.pid"
)

type Cmd string

const (
	CmdDiscover Cmd = "DISCOVER"
	CmdStop     Cmd = "STOP"
)

type CmdRequest struct {
	Cmd Cmd `json:"cmd"`
}

type CmdResponse struct {
	Peers []Peer `json:"peers,omitempty"`
	Error string `json:"error,omitempty"`
}

type Daemon struct {
	node *Node
}

func NewDaemon(port int) *Daemon {
	return &Daemon{node: NewNode(port)}
}

func (d *Daemon) Start() error {
	pid := os.Getpid()
	if err := os.WriteFile(pidPath, fmt.Appendf(nil, "%d", pid), 0644); err != nil {
		return fmt.Errorf("failed to write pid: %w", err)
	}

	if err := d.node.Start(); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	os.Remove(sockPath)
	ln, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("failed to listen on unix socket: %w", err)
	}

	log.Printf("[daemon] PID %d listening on %s", pid, sockPath)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[daemon] accept error:", err)
			continue
		}
		go d.handleConn(conn)
	}
}

func (d *Daemon) handleConn(conn net.Conn) {
	defer conn.Close()

	var req CmdRequest
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		log.Println("[daemon] bad request:", err)
		return
	}

	switch req.Cmd {
	case CmdDiscover:
		peers := d.node.Discover()
		json.NewEncoder(conn).Encode(CmdResponse{Peers: peers})

	case CmdStop:
		log.Println("[daemon] stop requested")
		json.NewEncoder(conn).Encode(CmdResponse{})
		os.Remove(sockPath)
		os.Remove(pidPath)
		os.Exit(0)
	}
}
