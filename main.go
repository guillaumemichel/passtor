package passtor

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

	p := NewPasstor(*name, *addr, *verbose)
	go p.ListenToPasstors()

	p.JoinDHT(ParsePeers(*peers))

	// keep the program active until ctrl+c is pressed
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.Printer.Print("", V0)
	os.Exit(0)
}
