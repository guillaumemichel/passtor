package passtor

import (
	"fmt"
	"math"
	"net"
	"sort"
	"sync"
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
		success = success || p.Ping(peer, MAXRETRIES)
	}
	if !success || len(peers) == 0 {
		if !success {
			p.Printer.WPrint("could not join the DHT, bootstrap peers do "+
				"not answer :(", V1)
		}
		p.Printer.Print("Creating new DHT", V1)
	} else {
		p.LookupReq(&p.NodeID)
		p.Printer.Print("Joined successfully the DHT", V1)
	}
}

// Ping a remote node, tries at most the given number of times
func (p *Passtor) Ping(peer net.UDPAddr, retries int) bool {
	b := true
	msg := Message{Ping: &b}
	// send it, this function returns once the reply is received
	return p.SendMessage(msg, peer, MAXRETRIES) != nil
}

// GetBucketID get the bucket identifier in which val belongs
func (p *Passtor) GetBucketID(val *Hash) uint16 {
	if p.NodeID.Compare(*val) == 0 {
		return 0
	}
	dist := p.NodeID.XOR(*val)

	// find in which bucket addr belongs
	var bucket uint16
	done := false
	for _, b := range dist {
		for i := BYTELENGTH - 1; i < BYTELENGTH; i-- {
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
	return bucket
}

// AddPeerToBucket check if a peer should be added to the DHT, and if yes,
// add it to the appropriate bucket
func (p *Passtor) AddPeerToBucket(addr NodeAddr) {
	bucket := p.GetBucketID(&addr.NodeID)
	// if bucket does not exist yet, create it
	if _, ok := p.Buckets[bucket]; !ok {
		p.Buckets[bucket] = &Bucket{
			Size:  0,
			Tail:  nil,
			Head:  nil,
			Mutex: &sync.Mutex{},
		}
	}
	// get the coresponding bucket
	b := p.Buckets[bucket]
	if el := b.Find(&addr); el != nil {
		// element already in bucket, moving it to head
		b.MoveToHead(el)
	} else if b.Size < DHTK {
		// element not in bucket && bucket not full, adding el to bucket
		b.Insert(&addr)
		p.Printer.Print(fmt.Sprint("Added ", addr, " to bucket ", bucket), V2)
	} else {
		// if the tail of the bucket does not reply the ping request, replace
		// it by the new address
		if !p.Ping(b.Tail.NodeAddr.Addr, MINRETRIES) {
			b.Replace(b.Tail, &addr)
		}
	}
}

// LookupReq lookup a hash
func (p *Passtor) LookupReq(hash *Hash) []NodeAddr {
	// set initial lookup peers
	initial := p.GetKCloser(hash)
	statuses := make([]*LookupStatus, len(initial))
	for i, na := range initial {
		statuses[i] = NewLookupStatus(na)
	}
	// TODO concurrency with concurrency parameter ALPHA
	found := true
	for found {
		found = false
		// if all values are already looked up, found = false -> exit loop
		for _, s := range statuses {
			if !s.Tested {
				// lookup at that node
				found = true
				msg := Message{LookupReq: hash}
				reply := p.SendMessage(msg, s.NodeAddr.Addr, MINRETRIES)
				// now we tested this node
				s.Tested = true
				if reply == nil {
					s.Failed = true
				} else {
					// update statuses with new results
					for _, n := range *reply.LookupRep {
						// if peer not in statuses yet, insert it
						alreadyIn := false
						for _, s := range statuses {
							if s.NodeAddr.NodeID.Compare(n.NodeID) == 0 {
								// already in list
								alreadyIn = true
								break
							}
						}
						if !alreadyIn && p.NodeID.Compare(n.NodeID) != 0 {
							// if not in list, add it
							statuses = append(statuses, NewLookupStatus(n))
						}
					}
				}
			}
		}
	}
	// get the K closest nodes

	// sort the array with id closest to target first
	sort.Slice(statuses, func(i, j int) bool {
		// dist(i, target) < dist(j, target)
		return statuses[i].NodeAddr.NodeID.XOR(*hash).Compare(
			statuses[j].NodeAddr.NodeID.XOR(*hash)) < 0
	})
	// number of elements to return
	n := DHTK
	if len(statuses) < DHTK {
		n = len(statuses)
	}
	results := make([]NodeAddr, n)
	for i := range results {
		results[i] = statuses[i].NodeAddr
	}
	return results
}

// LookupRep handles a lookup request and reply to it
func (p *Passtor) LookupRep(req *Message) {
	KCloser := p.GetKCloser(req.LookupReq)

	req.LookupReq = nil
	req.Reply = true
	req.LookupRep = &KCloser
	p.SendMessage(*req, req.Sender.Addr, MINRETRIES)
}

// GetKCloser get the K closer nodes to given hash
func (p *Passtor) GetKCloser(h *Hash) []NodeAddr {
	list := make([]NodeAddr, 0)
	var bID uint16 = math.MaxUint16

	if p.NodeID.Compare(*h) != 0 {
		bID = p.GetBucketID(h)

		b, ok := p.Buckets[bID]
		if ok {
			// bucket exists append all addresses in list
			list = append(list, b.GetList()...)
		}
	}
	if len(list) < DHTK {
		// less than k elements in corresponding bucket
		// complete with random elements
		for id, bucket := range p.Buckets {
			if id != bID {
				for _, h := range bucket.GetList() {
					list = append(list, h)
					if len(list) >= DHTK {
						break
					}
				}
			}
			if len(list) >= DHTK {
				break
			}
		}
	}
	return list
}

// Insert a new NodeAddress in the Bucket, called only if bucket not full
func (b *Bucket) Insert(nodeAddr *NodeAddr) {
	b.Mutex.Lock()
	if b.Size >= DHTK {
		fmt.Println("Warning: cannot insert node to full bucket!")
		b.Mutex.Unlock()
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
	b.Mutex.Unlock()
}

// Find and return the element corresponding to the given address, returns nil
// if not found
func (b *Bucket) Find(nodeAddr *NodeAddr) *BucketElement {
	b.Mutex.Lock()
	el := b.Tail
	for el != nil {
		if el.NodeAddr.NodeID == nodeAddr.NodeID {
			b.Mutex.Unlock()
			return el
		}
		el = el.Next
	}
	b.Mutex.Unlock()
	return nil
}

// Replace the a node address in the list by a new one
func (b *Bucket) Replace(old *BucketElement, new *NodeAddr) {
	b.Mutex.Lock()
	old.NodeAddr = new
	b.Mutex.Unlock()
}

// MoveToHead moves an element of the list to the head
func (b *Bucket) MoveToHead(el *BucketElement) {
	b.Mutex.Lock()
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
	b.Mutex.Unlock()
}

// GetList returns the list of node addresses in the given bucket
func (b *Bucket) GetList() []NodeAddr {
	b.Mutex.Lock()
	list := make([]NodeAddr, b.Size)
	el := b.Head
	for i := 0; el != nil; i++ {
		list[i] = *el.NodeAddr
		el = el.Prev
	}
	b.Mutex.Unlock()
	return list
}
