package passtor

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"go.dedis.ch/protobuf"
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

// ListenToPasstors listen on the udp connection used to communicate with other
// passtors, and distribute received messages to HandleMessage()
func (p *Passtor) ListenToPasstors() {
	buf := make([]byte, BUFFERSIZE)

	for {
		// read new message
		m, _, err := p.PConn.ReadFromUDP(buf)
		checkErr(err)

		// copy the receive buffer to avoid that it is modified while being used
		tmp := make([]byte, m)
		copy(tmp, buf[:m])

		go p.HandleMessage(tmp)
	}
}

func (p *Passtor) ListenToClients() {

	server, err := net.Listen("tcp", ":8080")
	accounts := make(Accounts)
	if err != nil {
		fmt.Println("Error while starting TCP server")
		return
	}
	defer server.Close()

	for {
		conn, _ := server.Accept()

		packetBytes := make([]byte, TCPMAXPACKETSIZE)
		_, err := conn.Read(packetBytes)
		if err != nil {
			println("Unable to read packet from TCP connection")
		}

		var message ClientMessage
		protobuf.Decode(packetBytes, &message)
		response := p.HandleClientMessage(accounts, message)
		responseBytes, err := protobuf.Encode(response)
		if err != nil {
			fmt.Println("Error while parsing response to be sent to client")
		}
		conn.Write(responseBytes)
	}
}
