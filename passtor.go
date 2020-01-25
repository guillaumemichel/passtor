package passtor

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

// NewPasstor creates and return a new Passtor instance
func NewPasstor(name, addr string, verbose int) Passtor {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkErr(err)
	pConn, err := net.ListenUDP("udp4", udpAddr)
	checkErr(err)

	// create the passtor printers
	printer := Printer{
		Verbose:    verbose,
		Printer:    log.New(os.Stdout, "", 0),
		ErrPrinter: log.New(os.Stderr, "", 0),
	}

	// create the message counter used to associate reply with request
	var counter uint64
	c := MessageCounter{
		IDCounter:  &counter,
		Mutex:      &sync.Mutex{},
		PendingMsg: make(map[uint64]*chan Message),
	}

	// create the passtor instance
	p := Passtor{
		Name:     name,
		Messages: &c,
		PConn:    pConn,
		Printer:  printer,
		Buckets:  make(map[uint16]*Bucket),
		Accounts: make(map[Hash]*AccountInfo),
	}
	// set the passtor identifier
	p.SetIdentity()
	// set self address
	p.Addr = NodeAddr{Addr: *udpAddr, NodeID: p.NodeID}
	// add self to routing table
	p.AddPeerToBucket(p.Addr)

	p.Printer.Print(fmt.Sprint("NodeID: ", p.NodeID), V3)
	return p
}

// GetMessageID get the next message ID, ids starting at 1
func (c MessageCounter) GetMessageID() uint64 {
	c.Mutex.Lock()
	*c.IDCounter++
	id := *c.IDCounter
	c.Mutex.Unlock()
	return id
}
