package main

import (
	"flag"
	"fmt"
	"gitlab.gnugen.ch/gmichel/passtor"
	"go.dedis.ch/protobuf"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handle(message passtor.ClientMessage) *passtor.ServerResponse {
	// TODO implement store and retrieve data logic
	return nil
}

func listenToClients() {

	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error while starting TCP server")
		return
	}
	conn, _ := server.Accept()

	go func() {
		for {
			packetBytes := make([]byte, passtor.TCPMAXPACKETSIZE)
			_, err := conn.Read(packetBytes)
			if err != nil {
				println("Unable to read packet from TCP connection")
			}

			var message passtor.ClientMessage
			protobuf.Decode(packetBytes, &message)
			response := handle(message)
			responseBytes, err := protobuf.Encode(response)
			if err != nil {
				fmt.Println("Error while parsing response to be sent to client")
			}
			conn.Write(responseBytes)
		}
	}()
}

func main() {

	name := flag.String("name", "", "name of the Passtor instance")
	addr := flag.String("addr", "127.0.0.1:5000", "address used to communicate "+
		"other passtors instances")
	peers := flag.String("peers", "", "bootstrap peer addresses")
	verbose := flag.Int("v", 1, "verbose mode")

	flag.Parse()
	// help message
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	p := passtor.NewPasstor(*name, *addr, *verbose)
	go p.ListenToPasstors()

	p.JoinDHT(passtor.ParsePeers(*peers))

	// keep the program active until ctrl+c is pressed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.Printer.Print("", passtor.V0)
	os.Exit(0)
}
