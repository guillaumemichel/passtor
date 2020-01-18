package main

import (
	"fmt"
	"log"
	"os"
)

// Printer of the passtor
var printer = log.New(os.Stdout, "", 0)

// Printer of the passtor
var errPrinter = log.New(os.Stderr, "", 0)

func checkErrMsg(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: "+msg+"\n")
		os.Exit(1)
	}
}

// Print a message to stdout
func Print(str string, V int) {
	if V <= VERBOSE {
		printer.Println(str)
	}
}

// WPrint prints warnings
func WPrint(str string) {
	printer.Println("Warning: " + str)
}

// PrintErr prints error message to stderr
func PrintErr(str string) {
	errPrinter.Println(str)
}
