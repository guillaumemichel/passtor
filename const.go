package main

import (
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

	// V0 verbose level 0
	V0 = 0
	// V1 verbose level 1
	V1 = 1
	// V2 verbose level 2
	V2 = 2
)

// VERBOSE level
var VERBOSE = 1
