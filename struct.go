package passtor

import (
	"log"
	"net"
	"sync"
)

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

	Addr    NodeAddr           // address used to communicate with passtors
	Buckets map[uint16]*Bucket // k-buckets used in the DHT

	Printer Printer // passtor console printer
}

// Message structure defining messages exchanged between passtors
type Message struct {
	ID        uint64      // message ID
	Reply     bool        // message is a reply
	Sender    *NodeAddr   // sender identity
	Ping      *bool       // non nil if message is a ping message
	LookupReq *Hash       // value to lookup
	LookupRep *[]NodeAddr // lookup response
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

// LoginMetaData for the Login structure.
type LoginMetaData struct {
	ServiceNonce  Nonce
	UsernameNonce Nonce
	PasswordNonce Nonce
}

// Credentials for a given service.
type Credentials struct {
	Username EncryptedData
	Password EncryptedData
}

// Login is a tuple of credentials and corresponding metadata to ensure validity.
type Login struct {
	ID          Hash
	Service     EncryptedData
	Credentials Credentials
	MetaData    LoginMetaData
}

// KeysClient used only client side to store the keys used to sign or en/de-crypt data.
type KeysClient struct {
	PublicKey    PublicKey
	PrivateKey   PrivateKey
	SymmetricKey SymmetricKey
}

// Keys used to encrypt, or sign data.
type Keys struct {
	PublicKey      PublicKey
	PrivateKeySeed EncryptedData
	SymmetricKey   EncryptedData
}

// AccountMetaData for the Account structure.
type AccountMetaData struct {
	PrivateKeySeedNonce Nonce
	SymmetricKeyNonce   Nonce
}

// Account used only client side to store info about the current user.
type AccountClient struct {
	ID string
	Keys KeysClient
}

// Account groups everything that has been stored by a single user.
type Account struct {
	ID        Hash
	Keys      Keys
	Version   uint
	Data      map[Hash]Login
	MetaData  AccountMetaData
	Signature Signature
}

// Accounts is the collection of all created accounts.
type Accounts map[Hash]Account
