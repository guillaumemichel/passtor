package passtor

import (
	"crypto/ed25519"
	"crypto/sha256"
	"golang.org/x/crypto/chacha20poly1305"
	"math"
	"time"
)

const (
	// DHTK K parameter for DHT K-buckets
	DHTK = 5
	// ALPHA DHT concurrency parameter
	ALPHA = 2
	// TIMEOUT value when waiting for an answer
	TIMEOUT = 2 * time.Second
	// MINRETRIES min number of attemps before giving up reaching an host
	MINRETRIES = 1
	// MAXRETRIES max number of attemps before giving up reaching an host
	MAXRETRIES = 4
	// BUFFERSIZE size of the udp connection read buffer
	BUFFERSIZE = 8192
	// SHASIZE size of SHA256 hash in byte
	SHASIZE = sha256.Size // TODO: change to SHA-3, change in crypto as well
	// BYTELENGTH number of bits in a byte
	BYTELENGTH uint16 = 8

	// SYMMKEYSIZE size in bytes for a symmetric key
	SYMMKEYSIZE = chacha20poly1305.KeySize
	// NONCESIZE size in bytes for a nonce
	NONCESIZE = chacha20poly1305.NonceSizeX
	// SIGNATURESIZE size in bytes for a signature
	SIGNATURESIZE = ed25519.SignatureSize

	// V0 verbose level 0 (no output)
	V0 = 0
	// V1 verbose level 1 (normal output)
	V1 = 1
	// V2 verbose level 2 (mode verbose)
	V2 = 2
	// V3 verbose level 3 (mode verbose++)
	V3 = 3
)

// MAXDISTANCE maximum distance between two hashes
var MAXDISTANCE Hash

func init() {
	// set MAXDISTANCE
	b := byte(math.MaxUint8)
	for i := range MAXDISTANCE {
		MAXDISTANCE[i] = b
	}
}
