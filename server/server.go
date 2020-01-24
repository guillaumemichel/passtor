package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"../../passtor"
	"go.dedis.ch/protobuf"
)

func handle(accounts passtor.Accounts, message passtor.ClientMessage) *passtor.ServerResponse {

	// Update or store new account
	if message.Push != nil {
		err := accounts.Store(*message.Push)
		if err == nil {
			return &passtor.ServerResponse{
				Status: "ok",
			}
		}

		msg := err.Error()
		return &passtor.ServerResponse{
			Status: "error",
			Debug:  &msg,
		}
	}

	// Retrieve account
	if message.Pull != nil {
		account, exists := accounts[*message.Pull]
		if exists {
			return &passtor.ServerResponse{
				Status: "ok",
				Data:   &account,
			}
		}
	}

	return nil

}

func listenToClients() {

	server, err := net.Listen("tcp", ":8080")
	accounts := make(passtor.Accounts)
	if err != nil {
		fmt.Println("Error while starting TCP server")
		return
	}
	defer server.Close()

	for {
		conn, _ := server.Accept()

		packetBytes := make([]byte, passtor.TCPMAXPACKETSIZE)
		_, err := conn.Read(packetBytes)
		if err != nil {
			println("Unable to read packet from TCP connection")
		}

		var message passtor.ClientMessage
		protobuf.Decode(packetBytes, &message)
		response := handle(accounts, message)
		responseBytes, err := protobuf.Encode(response)
		if err != nil {
			fmt.Println("Error while parsing response to be sent to client")
		}
		conn.Write(responseBytes)
	}
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
	go listenToClients()

	p.JoinDHT(passtor.ParsePeers(*peers))

	// keep the program active until ctrl+c is pressed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.Printer.Print("", passtor.V0)
	os.Exit(0)
}
