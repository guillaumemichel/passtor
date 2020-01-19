package passtor

import (
	"log"
	"net"
	"sync"
)

// Hash format of the sha256 hash function
type Hash [SHASIZE]byte

// NodeAddr node address entry in the k-bucket, node udp ip and port, and nodeID
type NodeAddr struct {
	Addr   net.UDPAddr // udp address (ip + port) of the node
	NodeID Hash        // nodeID of that node
}

// MessageCounter structure containing message indexing tools
type MessageCounter struct {
	Mutex *sync.Mutex // mutex of the structure

	IDCounter  uint64                   // current message ID
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
	Name   string // name of the passtor instance
	NodeID Hash   // hash of the name of the passtor, node identifier

	PConn *net.UDPConn // udp socket to communicate with other passtors
	CConn *net.UDPConn // udp socket to communicate with clients

	Messages MessageCounter // handles message id and pending messages

	Addr    NodeAddr         // address used to communicate with passtors
	Buckets map[uint]*Bucket // k-buckets used in the DHT

	Printer Printer // passtor console printer
}

// Message structure defining messages exchanged between passtors
type Message struct {
	ID        uint64    // message ID
	Reply     bool      // message is a reply
	Sender    *NodeAddr // sender identity
	Bootstrap *bool     // non nil if message is a bootstrap message
}

// Bucket structure representing Kademlia k-buckets
type Bucket struct {
	Head *BucketElement
	Tail *BucketElement
	Size uint
}

// BucketElement represent individual elements of the k-buckets
type BucketElement struct {
	NodeAddr *NodeAddr
	Next     *BucketElement
	Prev     *BucketElement
}
