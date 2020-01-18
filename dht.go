package main

import "net"

// JoinDHT passtor join the DHT
// connect to bootstrap peers given as argument
// lookup for self to fill k-buckets
func (p Passtor) JoinDHT(peers []net.UDPAddr) {

}

func (p Passtor) BootstrapPeer(peer net.UDPAddr) {
	msg := Message{Bootstrap: &true}
	p.SendMessage(msg, peer)
}

func (p Passtor) LookupNode(nodeID Hash) {

}
