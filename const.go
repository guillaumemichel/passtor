package passtor

import (
	"fmt"
	"math"
	"os"
	"time"
)

const (
	// DHTK K parameter for DHT K-buckets
	DHTK = 5
	// ALPHA DHT concurrency parameter
	ALPHA = 2
	// REPL replication factor
	REPL = 3
	// NREQ minimal number of response after Fetch
	NREQ = 3
	// THRESHOLD of answers before returning
	THRESHOLD = 0.333
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

	// PASSPHRASELENGHT default length in words for a passphrase
	PASSPHRASELENGHT = 8
	// PASSPHRASESEP default word seperator in a passphrase
	PASSPHRASESEP = "."

	// V0 verbose level 0 (no output)
	V0 = 0
	// V1 verbose level 1 (normal output)
	V1 = 1
	// V2 verbose level 2 (mode verbose)
	V2 = 2
	// V3 verbose level 3 (mode verbose++)
	V3 = 3
	// TCPMAXPACKETSIZE is the largest size in bytes of a TCP packet
	TCPMAXPACKETSIZE = 65535
	// REPUBLISHINTERVAL average time interval between republish
	REPUBLISHINTERVAL = 5 * time.Minute
)

// NOERROR string
var NOERROR = ""

// MAXDISTANCE maximum distance between two hashes
var MAXDISTANCE Hash

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
