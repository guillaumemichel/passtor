package passtor

import (
	"fmt"
	"math"
	"os"
	"time"

	"gitlab.gnugen.ch/gmichel/passtor/crypto"
)

const (
	// DHTK K parameter for DHT K-buckets
	DHTK = 5
	// ALPHA DHT concurrency parameter
	ALPHA = 2
	// REPL replication factor
	REPL = 3
	// TIMEOUT value when waiting for an answer
	TIMEOUT = 2 * time.Second
	// MINRETRIES min number of attemps before giving up reaching an host
	MINRETRIES = 1
	// MAXRETRIES max number of attemps before giving up reaching an host
	MAXRETRIES = 4
	// BUFFERSIZE size of the udp connection read buffer
	BUFFERSIZE = 8192
	// BYTELENGTH number of bits in a byte
	BYTELENGTH uint16 = 8

	// V0 verbose level 0 (no output)
	V0 = 0
	// V1 verbose level 1 (normal output)
	V1 = 1
	// V2 verbose level 2 (mode verbose)
	V2 = 2
	// V3 verbose level 3 (mode verbose++)
	V3 = 3

	// NOERROR string
	NOERROR = ""
	// REPUBLISHINTERVAL average time interval between republish
	REPUBLISHINTERVAL = 5 * time.Minute
)

// MAXDISTANCE maximum distance between two hashes
var MAXDISTANCE crypto.Hash

func init() {
	// set MAXDISTANCE
	b := byte(math.MaxUint8)
	for i := range MAXDISTANCE {
		MAXDISTANCE[i] = b
	}
	if REPL > DHTK {
		fmt.Println("Replication factor can't be larger than K")
		os.Exit(1)
	}
}
