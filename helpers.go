package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkErrMsg(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: "+msg+"\n")
		os.Exit(1)
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
