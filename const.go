package passtor

import (
	"crypto/sha256"
	"time"
)

const (
	// DHTK K parameter for DHT K-buckets
	DHTK = 5
	// TIMEOUT value when waiting for an answer
	TIMEOUT = 2 * time.Second
	// MAXRETRIES number of attemps before giving up reaching an host
	MAXRETRIES = 4
	// BUFFERSIZE size of the udp connection read buffer
	BUFFERSIZE = 8192
	// SHASIZE size of SHA256 hash in byte
	SHASIZE = sha256.Size
	// BYTELENGTH number of bits in a byte
	BYTELENGTH = 8

	// V0 verbose level 0
	V0 = 0
	// V1 verbose level 1
	V1 = 1
	// V2 verbose level 2
	V2 = 2
)
