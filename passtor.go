package main

import "net"

// NewPasstor creates and return a new Passtor instance
func NewPasstor(name, addr string) Passtor {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkErr(err)
	p := Passtor{Name: name}
	p.SetIdentity()
	p.Addr = NodeAddr{Addr: udpAddr, NodeID: p.NodeID}
	return p
}
