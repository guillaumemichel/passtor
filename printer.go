package passtor

import (
	"fmt"
	"os"
)

func checkErrMsg(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: "+msg+"\n")
		os.Exit(1)
	}
}

// Print a message to stdout
func (p *Printer) Print(str string, V int) {
	if V <= p.Verbose {
		p.Printer.Println(str)
	}
}

// WPrint prints warnings
func (p *Printer) WPrint(str string, V int) {
	if V <= p.Verbose {
		p.Printer.Println("Warning: " + str)
	}
}

// PrintErr prints error message to stderr
func (p *Printer) PrintErr(str string) {
	p.ErrPrinter.Println(str)
}
