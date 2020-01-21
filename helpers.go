package passtor

import (
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

// NewLookupStatus returns new lookup status structure for given nodeaddr
func NewLookupStatus(nodeAddr NodeAddr) *LookupStatus {
	return &LookupStatus{
		NodeAddr: nodeAddr,
		Failed:   false,
		Tested:   false,
	}
}
