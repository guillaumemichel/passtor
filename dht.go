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

	var bucket uint
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
	if _, ok := p.Buckets[bucket]; !ok {
		p.Buckets[bucket] = &Bucket{
			Size: 0,
			Tail: nil,
			Head: nil,
		}
	}
	b := p.Buckets[bucket]
	if b.Size < DHTK {
		b.Insert(&addr)
		p.Printer.Print(fmt.Sprint("Added ", addr, " to bucket ", bucket), V2)
	} else {
		p.Printer.Print("Bucket full", V2)
	}
}

// CheckPeersAlive check if DHT peers are alive, and remove them from the
// k-bucket if they don't seem to be alive
func (p *Passtor) CheckPeersAlive() {

}

// LookupNode lookup a node by its ID
func (p *Passtor) LookupNode(nodeID Hash) {

}

// Insert a new NodeAddress in the Bucket, called only if bucket not full
func (b *Bucket) Insert(nodeAddr *NodeAddr) {
	if b.Size >= DHTK {
		fmt.Println("Warning: cannot insert node to full bucket!")
		return
	}
	new := BucketElement{
		NodeAddr: nodeAddr,
		Next:     b.Tail,
		Prev:     nil,
	}
	if b.Size == 0 {
		b.Head = &new
	} else {
		b.Tail.Prev = &new
	}
	b.Tail = &new
	b.Size++
}

// Find and return the element corresponding to the given address, returns nil
// if not found
func (b *Bucket) Find(nodeAddr *NodeAddr) *BucketElement {
	el := b.Tail
	for el != nil {
		if el.NodeAddr.NodeID == nodeAddr.NodeID {
			return el
		}
		el = el.Next
	}
	return nil
}

// Replace the a node address in the list by a new one
func (b *Bucket) Replace(old *BucketElement, new *NodeAddr) {
	if b.Size == 0 {
		fmt.Println("Warning: no element to replace in bucket")
	}
	old.NodeAddr = new
}

// MoveToHead moves an element of the list to the head
func (b *Bucket) MoveToHead(el *BucketElement) {
	if el != b.Head {
		el.Next.Prev = el.Prev

		if el != b.Tail {
			el.Prev.Next = el.Next
		} else {
			b.Tail = el.Next
		}
		b.Head.Next = el
		el.Prev = b.Head
		el.Next = nil
		b.Head = el
	}
}
