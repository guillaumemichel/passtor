package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	name := flag.String("name", "", "name of the Passtor instance")
	addr := flag.String("addr", "127.0.0.1", "address used to communicate "+
		"other passtors instances")
	peers := flag.String("peers", "", "bootstrap peer addresses")

	flag.Parse()
	// help message
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	p := NewPasstor(*name, *addr)
	ParsePeers(*peers)

	fmt.Println(p)
}
