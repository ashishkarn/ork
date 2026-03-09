package main

import "sync"

type Peer struct {
	ID   string
	Addr string
}

type PeerMap struct {
	mu    sync.RWMutex
	peers map[string]Peer
}

func NewPeerMap() *PeerMap {
	return &PeerMap{peers: make(map[string]Peer)}
}

func (p *PeerMap) Add(peer Peer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers[peer.ID] = peer
}

func (p *PeerMap) All() []Peer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	list := make([]Peer, 0, len(p.peers))
	for _, peer := range p.peers {
		list = append(list, peer)
	}
	return list
}
