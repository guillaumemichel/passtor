package passtor

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParsePeers parse peer list in string format to udp addresses
func ParsePeers(peerList string) []net.UDPAddr {

	addresses := make([]net.UDPAddr, 0)

	if peerList == "" {
		return addresses
	}
	// split up the different addresses
	peers := strings.Split(peerList, ",")

	// parse the addresses and add them to the slice
	for _, p := range peers {
		udpAddr, err := net.ResolveUDPAddr("udp4", p)
		checkErrMsg(err, "invalid address \""+p+"\"")
		addresses = append(addresses, *udpAddr)
	}

	return addresses
}

// Timeout creates a clock that writes to the returned channel after the
// time value given as argument
func Timeout(timeout time.Duration) *chan bool {
	c := make(chan bool)
	go func() {
		time.Sleep(timeout)
		c <- true
	}()
	return &c
}

// XOR function computing the XOR distance between two hashes
func (hash0 Hash) XOR(hash1 Hash) Hash {
	res := Hash{}
	for i := range hash0[:] {
		res[i] = hash0[i] ^ hash1[i]
	}
	return res
}

// Compare two hashes, returns 1 if first hash smaller than the second, -1 if
// the second is smaller than the first, and 0 if they are equal
func (hash0 Hash) Compare(hash1 Hash) int {
	if len(hash0) != len(hash1) {
		fmt.Println("Cannot compare hashes of different sizes")
		return 0
	}
	for i := 0; i < len(hash0); i++ {
		if hash0[i] < hash1[i] {
			return -1
		} else if hash0[i] > hash1[i] {
			return 1
		}
	}
	return 0
}

// NewLookupStatus returns new lookup status structure for given nodeaddr
func NewLookupStatus(nodeAddr NodeAddr) *LookupStatus {
	return &LookupStatus{
		NodeAddr: nodeAddr,
		Failed:   false,
		Tested:   false,
	}
}
