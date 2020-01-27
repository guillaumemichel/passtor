package passtor

import (
	"fmt"
	"os"
)

func checkErrMsg(err error, msg string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: "+msg+"\n")
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

// PrintBuckets print all bucket with their state
func (p *Passtor) PrintBuckets() {
	str := "----------\nPrinting Bucket\n----------\n"
	for i, b := range p.Buckets {
		str += fmt.Sprintln("Bucket", i, ":")
		list := b.GetList()
		for _, n := range list {
			str += fmt.Sprintln(p.NodeID.XOR(n.NodeID), n.NodeID, n.Addr)
		}
		str += "\n"
	}
	str += "----------"
	p.Printer.Print(str, V3)
}

// PrintStatuses print given lookup statuses
func PrintStatuses(h Hash, statuses []*LookupStatus) {
	str := "Printing lookup statuses:\n"
	for _, s := range statuses {
		str += fmt.Sprintln(*s)
		str += fmt.Sprintln(s.NodeAddr.NodeID.XOR(h).Hex())
	}
	fmt.Print(str)
}
