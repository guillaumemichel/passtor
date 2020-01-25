package passtor

import (
	"log"
	"net"
	"sync"

	"gitlab.gnugen.ch/gmichel/passtor/crypto"
)

// NodeAddr node address entry in the k-bucket, node udp ip and port, and nodeID
type NodeAddr struct {
	Addr   net.UDPAddr // udp address (ip + port) of the node
	NodeID crypto.Hash // nodeID of that node
}

// MessageCounter structure containing message indexing tools
type MessageCounter struct {
	Mutex *sync.Mutex // mutex of the structure

	IDCounter  *uint64                  // current message ID
	PendingMsg map[uint64]*chan Message // list of current pending messages
}

// Printer of the passtor, handles all prints to console
type Printer struct {
	Verbose    int
	Printer    *log.Logger
	ErrPrinter *log.Logger
}

// Passtor instance
type Passtor struct {
	Name   string      // name of the passtor instance
	NodeID crypto.Hash // hash of the name of the passtor, node identifier

	PConn *net.UDPConn // udp socket to communicate with other passtors
	CConn *net.UDPConn // udp socket to communicate with clients

	Messages *MessageCounter // handles message id and pending messages

	Addr    NodeAddr           // address used to communicate with passtors
	Buckets map[uint16]*Bucket // k-buckets used in the DHT

	Printer Printer // passtor console printer
}

// Message structure defining messages exchanged between passtors
type Message struct {
	ID            uint64       // message ID
	Reply         bool         // message is a reply
	Sender        *NodeAddr    // sender identity
	Ping          *bool        // non nil if message is a ping message
	LookupReq     *crypto.Hash // value to lookup
	LookupRep     *[]NodeAddr  // lookup response
	AllocationReq *AllocateMessage
	AllocationRep *string
	FetchReq      *crypto.Hash
	FetchRep      *Datastructure
}

// Bucket structure representing Kademlia k-buckets
type Bucket struct {
	Mutex *sync.Mutex
	Head  *BucketElement
	Tail  *BucketElement
	Size  uint
}

// BucketElement represent individual elements of the k-buckets
type BucketElement struct {
	NodeAddr *NodeAddr
	Next     *BucketElement
	Prev     *BucketElement
}

// LookupStatus type used by the lookup RPC
type LookupStatus struct {
	NodeAddr NodeAddr
	Tested   bool
	Failed   bool
}

// AllocateMessage message requesting a node to allocate a file
type AllocateMessage struct {
	Data  Datastructure
	Index uint32
	Repl  uint32
}

// Datastructure temporary structre of things that will be written somewhere
type Datastructure struct {
	MyData string
}
