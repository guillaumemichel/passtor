package passtor

import (
	"fmt"
	"net"
)

// JoinDHT passtor join the DHT
// connect to bootstrap peers given as argument
// lookup for self to fill k-buckets
func (p *Passtor) JoinDHT(peers []net.UDPAddr) {
	// check if at least a peer succeeded to contact a host
	// if no boostrap peer, create own DHT, always a success
	success := len(peers) == 0
	for _, peer := range peers {
		// if at least one peer succeed, the join is a success
		success = success || p.BootstrapPeer(peer)
	}
	if !success || len(peers) == 0 {
		if !success {
			p.Printer.WPrint("could not join the DHT, bootstrap peers do "+
				"not answer :(", V1)
		}
		p.Printer.Print("Creating new DHT", V1)
	} else {
		p.LookupNode(p.NodeID)
		p.Printer.Print("Joined successfully the DHT", V1)
	}
}

// BootstrapPeer send bootstrap message to the given peer
// return true if peer could be reached
func (p *Passtor) BootstrapPeer(peer net.UDPAddr) bool {
	// define bootstrap message
	b := true
	msg := Message{Bootstrap: &b}
	// send it, this function returns once the reply is received
	return p.SendMessage(msg, peer) != nil
}

// AddPeerToBucket check if a peer should be added to the DHT, and if yes,
// add it to the appropriate bucket
func (p *Passtor) AddPeerToBucket(addr NodeAddr) {
	dist := XOR(p.NodeID, addr.NodeID)

	bucket := 0
	done := false
	for _, b := range dist {
		for i := BYTELENGTH - 1; i >= 0; i-- {
			if b>>i > 0 {
				bucket += BYTELENGTH - i - 1
				done = true
				break
			}
		}
		if done {
			break
		}
		bucket += 8
	}
	fmt.Println(bucket)
	p.Printer.Print(fmt.Sprint("Adding", addr, "to k-bucket"), V2)
}

// CheckPeersAlive check if DHT peers are alive, and remove them from the
// k-bucket if they don't seem to be alive
func (p *Passtor) CheckPeersAlive() {

}

// LookupNode lookup a node by its ID
func (p *Passtor) LookupNode(nodeID Hash) {

}

// Insert a new NodeAddress in the Bucket, called only if bucket not full
func (b *Bucket) Insert(nodeAddr NodeAddr) {
	if b.Size >= DHTK {
		fmt.Println("Warning: cannot insert node to full bucket!")
		return
	}
	b.Size++
	new := BucketElement{
		NodeAddr: &nodeAddr,
		Next:     b.Tail,
	}
	b.Tail = &new
}

// Replace the tail of the list by the new Node Address
func (b *Bucket) Replace(nodeAddr NodeAddr) {
	new := BucketElement{
		NodeAddr: &nodeAddr,
		Next:     b.Tail.Next,
	}
	b.Tail = &new
}

// MoveTailToHead moves the tail element of the list to the head
func (b *Bucket) MoveTailToHead() {
	newHead := b.Tail
	b.Tail = b.Tail.Next
	newHead.Next = nil
	b.Head.Next = newHead
	b.Head = newHead
}
