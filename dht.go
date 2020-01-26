package passtor

import (
	"errors"
	"fmt"
	"math"
	"net"
	"sort"
	"sync"
	"time"
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
		return HASHSIZE
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
		p.Printer.Print(fmt.Sprint("Added ", addr, " to bucket ", bucket), V3)
	} else {
		// if the tail of the bucket does not reply the ping request, replace
		// it by the new address
		if !p.Ping(b.Tail.NodeAddr.Addr, MINRETRIES) {
			b.Replace(b.Tail, &addr)
		}
	}
	p.PrintBuckets()
}

// LookupReq lookup a hash
func (p *Passtor) LookupReq(hash *Hash) []NodeAddr {
	// set initial lookup peers
	initial := p.GetKCloser(hash)
	statuses := make([]*LookupStatus, len(initial))
	m := sync.Mutex{}
	for i, na := range initial {
		statuses[i] = NewLookupStatus(na)
		if statuses[i].NodeAddr.NodeID == p.Addr.NodeID {
			statuses[i].Tested = true
		}
	}
	wg := sync.WaitGroup{}

	for i := 0; i < ALPHA; i++ {
		// ALPHA parallel requests iterating on the statuses
		go func() {
			wg.Add(1)
			found := true
			for found {
				found = false
				// if all values are already looked up, found = false -> exit loop
				for _, s := range statuses {
					m.Lock()
					if !s.Tested {
						// lookup at that node
						s.Tested = true
						m.Unlock()
						found = true
						msg := Message{LookupReq: hash}
						reply := p.SendMessage(msg, s.NodeAddr.Addr, MINRETRIES)
						// now we tested this node
						if reply == nil {
							s.Failed = true
						} else {
							// update statuses with new results
							for _, n := range *reply.LookupRep {
								// if peer not in statuses yet, insert it
								alreadyIn := false
								m.Lock()
								for _, s := range statuses {
									if s.NodeAddr.NodeID.Compare(n.NodeID) == 0 {
										// already in list
										alreadyIn = true
										break
									}
								}
								m.Unlock()
								if !alreadyIn && p.NodeID.Compare(n.NodeID) != 0 {
									// if not in list, add it
									m.Lock()
									statuses = append(statuses, NewLookupStatus(n))
									m.Unlock()
								}
							}
						}
					} else {
						m.Unlock()
					}
				}
			}
			wg.Done()
		}()
	}
	// waiting for the ALPHA threads to finish
	wg.Wait()
	// get the K closest nodes

	// sort the array with id closest to target first
	sort.Slice(statuses, func(i, j int) bool {
		// (!fail[i] && fail[j]) ||
		//        (fail[i]==fail[j] && (dist(i, target) < dist(j, target)))
		return (!statuses[i].Failed && statuses[j].Failed) &&
			(statuses[i].Failed == statuses[j].Failed &&
				statuses[i].NodeAddr.NodeID.XOR(*hash).Compare(
					statuses[j].NodeAddr.NodeID.XOR(*hash)) < 0)
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
func (p *Passtor) LookupRep(req Message) {
	KCloser := p.GetKCloser(req.LookupReq)

	req.LookupReq = nil
	req.LookupRep = &KCloser
	p.SendMessage(req, req.Sender.Addr, MINRETRIES)
}

// HandleAllocation handles an allocation on the remote peer
func (p *Passtor) HandleAllocation(msg Message) {
	req := msg.AllocationReq

	msg.AllocationReq = nil
	err := p.Store(req.Account.ToAccount(), req.Repl)
	if err == nil {
		go p.Republish(p.Accounts[req.Account.ID])
		msg.AllocationRep = &NOERROR
	} else {
		errMsg := err.Error()
		msg.AllocationRep = &errMsg
	}

	p.SendMessage(msg, msg.Sender.Addr, MINRETRIES)
}

// AllocateToPeer allocate some data to a peer, returns true on success,
// false if cannot reach peer or error
func (p *Passtor) AllocateToPeer(id Hash, peer NodeAddr, index, repl uint32,
	data AccountNetwork) error {

	msg := Message{AllocationReq: &AccountMessage{
		Account: data,
		Repl:    repl,
	}}
	rep := p.SendMessage(msg, peer.Addr, MINRETRIES)
	if rep == nil {
		return errors.New("No response from " + peer.Addr.String())
	} else if rep.AllocationRep == nil {
		return errors.New("Invalid response from " + peer.Addr.String())
	} else if *rep.AllocationRep == NOERROR {
		return nil
	}
	return errors.New(*rep.AllocationRep)
}

// Allocate given data identified by the given id to the given replication
// factor appropriate peers
func (p *Passtor) Allocate(id Hash, repl uint32, data AccountNetwork) []NodeAddr {
	peers := p.LookupReq(&id)

	count := 0
	allocations := make([]NodeAddr, 0)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}

	limit := int(repl)
	if len(peers) < limit {
		limit = len(peers)
	}
	wg.Add(limit)

	for i := 0; i < limit; i++ {
		go func() {
			// while repl factor not match and not all peers tried
			m.Lock()
			for count < len(peers) {
				// take an element in the list and increase the counter
				peer := peers[count]
				count++
				index := uint32(len(allocations))
				m.Unlock()
				// allocation was a
				if err := p.AllocateToPeer(id, peer, index, repl, data); err == nil {
					m.Lock()
					allocations = append(allocations, peer)
					break
				} else {
					p.Printer.ErrPrinter.Println(err.Error())
				}
				m.Lock()
			}
			m.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	p.Printer.Print(fmt.Sprint("Allocated at:", allocations), V2)
	return allocations
}

// HandleFetch searches for requested file, send it if it finds it
func (p *Passtor) HandleFetch(msg Message) {
	p.Printer.Print("Got fetch req", V2)

	// if file is stored, return it
	if data, ok := p.Accounts[*msg.FetchReq]; ok {
		p.Printer.Print("I got it", V2)
		msg.FetchRep = &AccountMessage{
			Account: data.Account.ToAccountNetwork(),
			Repl:    data.Repl,
		}
	} else {
		p.Printer.Print("didn't found the data", V2)
	}
	// writing reply
	msg.FetchReq = nil

	// sending reply
	p.SendMessage(msg, msg.Sender.Addr, MINRETRIES)
}

// FetchDataFromPeer send fetch request to given peer, returns the reply of the
// remote host
func (p *Passtor) FetchDataFromPeer(h *Hash, peer NodeAddr) *Message {
	return p.SendMessage(Message{FetchReq: h}, peer.Addr, MINRETRIES)
}

// FetchData request associated with given hash from the DHT
func (p *Passtor) FetchData(h *Hash, threshold float64) *Account {
	peers := p.LookupReq(h)

	var min int
	var count int
	done := false
	replies := make([]Account, 0)
	var account Account
	m := sync.Mutex{}
	wg := sync.WaitGroup{}

	// TODO
	// waits for at least NREQ answers before calling MostRepresented(),
	// take most frequent repl
	p.Printer.Print("Fetching data...", V2)

	for i := 0; i < REPL; i++ {
		wg.Add(1)
		go func() {
			m.Lock()
			for !done && count < len(peers) {
				peer := peers[count]
				count++
				m.Unlock()
				if rep := p.FetchDataFromPeer(h, peer); rep != nil {
					if d := rep.FetchRep; d != nil {
						m.Lock()
						p.Printer.Print(fmt.Sprint(d), V2)
						if !done {
							p.Printer.Print("not done yet", V2)
							min = int(math.Ceil(threshold * float64(d.Repl)))
							p.Printer.Print(fmt.Sprint("min:", min), V2)
							p.Printer.Print(fmt.Sprint("cast min:", int(math.Ceil(threshold*float64(d.Repl)))), V2)
							p.Printer.Print(fmt.Sprint("ceil min:", math.Ceil(threshold*float64(d.Repl))), V2)
							p.Printer.Print(fmt.Sprint("float min:", threshold*float64(d.Repl)), V2)
							replies = append(replies, d.Account.ToAccount())
							if acc, ok := MostRepresented(replies, min); ok {
								p.Printer.Print("done now", V2)
								account = *acc
								done = true
								break
							}
							m.Unlock()
						} else {
							break
						}
					}
				}
				m.Lock()
			}
			m.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	if !done {
		acc, _ := MostRepresented(replies, min)
		p.Printer.Print(fmt.Sprint("min:", min), V2)
		p.Printer.Print(fmt.Sprint(acc), V2)
		return acc
	}
	p.Printer.Print(fmt.Sprint("Account:", account), V2)
	return &account
}

// Republish account information in the DHT, called periodically
func (p *Passtor) Republish(account *AccountInfo) {

	// gen random delay between 0 and 2*repl*REPUBLISHINTERVAL minutes
	delay := RandInt(int64(account.Repl) * int64(REPUBLISHINTERVAL) * 2 *
		int64(time.Minute/time.Second))
	// sleep for that time
	time.Sleep(time.Duration(delay) * time.Second)

	// check if version changed in the meantime, and republish only if it didn't
	republish := false
	p.Accounts[account.Account.ID].Mutex.Lock()
	if p.Accounts[account.Account.ID].Account.Version == account.Account.Version {
		// version did not change
		republish = true
	}
	p.Accounts[account.Account.ID].Mutex.Unlock()

	if republish {
		// remove account from the passtor
		p.Delete(account.Account.ID)

		// publish it
		addrs := p.Allocate(account.Account.ID, account.Repl,
			account.Account.ToAccountNetwork())
		// check if allocated to enough peers
		if len(addrs) != int(account.Repl) {
			p.Printer.WPrint("Couldn't reallocate to enough peers", V2)
		}
	}
}

// GetKCloser get the K closer nodes to given hash
func (p *Passtor) GetKCloser(h *Hash) []NodeAddr {

	if b, ok := p.Buckets[p.GetBucketID(h)]; ok && b.Size == DHTK {
		// bucket exists append all addresses in list
		return b.GetList()
	}

	list := make([]NodeAddr, 0)

	for _, b := range p.Buckets {
		list = append(list, b.GetList()...)
	}

	// sort the array with id closest to target first
	sort.Slice(list, func(i, j int) bool {
		// (i xor h) < (j xor h)
		return list[i].NodeID.XOR(*h).Compare(list[j].NodeID.XOR(*h)) < 0
	})

	size := DHTK
	if len(list) < DHTK {
		size = len(list)
	}

	return list[:size]
}

// Insert a new NodeAddress in the Bucket, called only if bucket not full
func (b *Bucket) Insert(nodeAddr *NodeAddr) {
	b.Mutex.Lock()
	if b.Size >= DHTK {
		fmt.Println("Warning: cannot insert node to full bucket!")
		b.Mutex.Unlock()
		return
	}
	newE := BucketElement{
		NodeAddr: nodeAddr,
		Next:     b.Tail,
		Prev:     nil,
	}
	if b.Size == 0 {
		b.Head = &newE
	} else {
		b.Tail.Prev = &newE
	}
	b.Tail = &newE
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
