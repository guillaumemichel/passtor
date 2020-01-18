package main

import (
	"net"
	"sync"
)

// NewPasstor creates and return a new Passtor instance
func NewPasstor(name, addr string) Passtor {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkErr(err)
	pConn, err := net.ListenUDP("udp4", udpAddr)
	checkErr(err)

	// create the message counter used to associate reply with request
	c := MessageCounter{
		IDCounter:  0,
		Mutex:      &sync.Mutex{},
		PendingMsg: make(map[uint64]*chan Message),
	}

	// create the passtor instance
	p := Passtor{
		Name:     name,
		Messages: c,
		PConn:    pConn,
	}
	// set the passtor identifier
	p.SetIdentity()
	p.Addr = NodeAddr{Addr: *udpAddr, NodeID: p.NodeID}
	return p
}

// GetMessageID get the next message ID, ids starting at 1
func (c MessageCounter) GetMessageID() uint64 {
	c.Mutex.Lock()
	c.IDCounter++
	id := c.IDCounter
	c.Mutex.Unlock()
	return id
}
