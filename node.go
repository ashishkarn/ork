package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Node struct {
	self    Peer
	peers   *PeerMap
	udpConn *net.UDPConn
}

func NewNode(port int) *Node {
	hostname, _ := os.Hostname()
	return &Node{
		self: Peer{
			ID:   hostname,
			Port: uint16(port),
		},
		peers: NewPeerMap(),
	}
}

func (n *Node) Start() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", n.self.Port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	n.udpConn = conn
	log.Printf("[node] UDP listening on %d as %s", n.self.Port, n.self.ID)
	go n.listenUDP()
	return nil
}

func (n *Node) listenUDP() {
	buf := make([]byte, 1024)
	for {
		sz, remoteAddr, err := n.udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("[udp] read error:", err)
			continue
		}

		msg, err := Decode(buf[:sz])
		if err != nil {
			log.Println("[udp] decode error:", err)
			continue
		}

		switch msg.Type {
		case MsgDiscover:
			if msg.NodeID == n.self.ID {
				continue
			}
			log.Printf("[udp] DISCOVER from %s (%s)", msg.NodeID, remoteAddr)
			n.sendAnnounce(remoteAddr)

		case MsgAnnounce:
			if msg.NodeID == n.self.ID {
				continue
			}
			peer := Peer{
				ID:   msg.NodeID,
				Addr: remoteAddr.IP.String(),
				Port: msg.Port,
			}
			n.peers.Add(peer)
			log.Printf("[udp] ANNOUNCE from %s @ %s:%d", peer.ID, peer.Addr, peer.Port)
		}
	}
}

func (n *Node) sendAnnounce(to *net.UDPAddr) {
	msg := &Message{
		Type:   MsgAnnounce,
		NodeID: n.self.ID,
		Port:   n.self.Port,
	}
	n.udpConn.WriteToUDP(msg.Encode(), to)
	log.Printf("[udp] ANNOUNCE sent to %s", to)
}

func (n *Node) Discover() []Peer {
	msg := &Message{
		Type:   MsgDiscover,
		NodeID: n.self.ID,
	}

	broadcast, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", n.self.Port))
	n.udpConn.WriteToUDP(msg.Encode(), broadcast)
	log.Printf("[discover] DISCOVER broadcast sent, collecting for 2s...")

	time.Sleep(2 * time.Second)
	return n.peers.All()
}
