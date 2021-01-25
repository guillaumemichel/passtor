package main

import (
	"fmt"
	"net"

	"github.com/guillaumemichel/passtor"
	"go.dedis.ch/protobuf"
)

// Contains returns true iff the given array is not empty and contains the given value
func Contains(array []string, value string) bool {

	if len(array) == 0 {
		return false
	}

	for _, str := range array {
		if str == value {
			return true
		}
	}

	return false
}

// AbortOnError immediately exits the program after printing the error message
func AbortOnError(err error, message string) {
	if err != nil {
		fmt.Println("Error", message)
		client.App.Stop()
	}
}

// FailWithError terminates the program after printing an error
func FailWithError(message string, debug *string) {
	fmt.Println(message)
	if debug != nil {
		fmt.Println(*debug)
	}
	client.App.Stop()
}

// Request queries the given node with the given message and returns the response from the server
func Request(message *passtor.ClientMessage, host string) *passtor.ServerResponse {

	conn, err := net.Dial("tcp", host)
	AbortOnError(err, "Unable to reach node")

	packet, err := protobuf.Encode(message)
	AbortOnError(err, "Could not encode message")

	_, err = conn.Write(packet)
	AbortOnError(err, "Could not send packet to node")

	reply := make([]byte, passtor.TCPMAXPACKETSIZE)
	n, err := conn.Read(reply)
	AbortOnError(err, "Could not read packet from node")

	var response passtor.ServerResponse
	err = protobuf.Decode(reply[:n], &response)
	AbortOnError(err, "Could not parse server response")

	return &response

}
