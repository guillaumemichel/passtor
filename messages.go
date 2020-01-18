package main

import (
	"net"

	"github.com/dedis/protobuf"
)

// SendMessage send the given message to the remote peer over udp
func (p Passtor) SendMessage(msg Message, dst net.UDPAddr) {
	msg.Sender = &p.Addr
	bytes, err := protobuf.Encode(&msg)
	checkErr(err)
	udpConn, err := net.DialUDP("udp4", nil, dst)
	checkErr(err)
	_, err := udpConn.Write(bytes)
	checkErr(err)
}
