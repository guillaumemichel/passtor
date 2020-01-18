package main

import (
	"crypto/sha256"
	"net"
	"sync"
)

// Hash format of the sha256 hash function
type Hash [sha256.Size]byte

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

// Passtor instance
type Passtor struct {
	Name   string // name of the passtor instance
	NodeID Hash   // hash of the name of the passtor, node identifier

	PConn *net.UDPConn // udp socket to communicate with other passtors
	CConn *net.UDPConn // udp socket to communicate with clients

	Messages MessageCounter // handles message id and pending messages

	Addr    NodeAddr     // address used to communicate with other passtors
	Buckets [][]NodeAddr // k-buckets used in the DHT

}

// Message structure defining messages exchanged between passtors
type Message struct {
	ID        uint64    // message ID
	Reply     bool      // message is a reply
	Sender    *NodeAddr // sender identity
	Bootstrap *bool     // non nil if message is a bootstrap message
}
