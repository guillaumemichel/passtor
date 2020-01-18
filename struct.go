package main

import (
	"crypto/sha256"
	"net"
)

// Hash format of the sha256 hash function
type Hash [sha256.Size]byte

// NodeAddr node address entry in the k-bucket, node udp ip and port, and nodeID
type NodeAddr struct {
	Addr   net.UDPAddr
	NodeID Hash
}

// Passtor instance
type Passtor struct {
	Name    string       // name of the passtor instance
	NodeID  Hash         // hash of the name of the passtor, node identifier
	Addr    NodeAddr     // address used to communicate with other passtors
	Buckets [][]NodeAddr // k-buckets used in the DHT
}

// Message structure defining messages exchanged between passtors
type Message struct {
	Sender    *NodeAddr
	Bootstrap *bool
}
