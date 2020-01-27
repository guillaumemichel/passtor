package passtor

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// HASHSIZE size of a hash in byte
const HASHSIZE = 64

// Hash is a flexible type to handle hashes
type Hash [HASHSIZE]byte

// H hashes the given bytes value
func H(data []byte) Hash {
	return sha3.Sum512(data)
}

// XOR function computing the XOR distance between two hashes
func (hash0 Hash) XOR(hash1 Hash) Hash {
	res := Hash{}
	for i := range hash0[:] {
		res[i] = hash0[i] ^ hash1[i]
	}
	return res
}

// Compare two hashes, returns 1 if first hash smaller than the second, -1 if
// the second is smaller than the first, and 0 if they are equal
func (hash0 Hash) Compare(hash1 Hash) int {
	for i := 0; i < len(hash0); i++ {
		if hash0[i] < hash1[i] {
			return -1
		} else if hash0[i] > hash1[i] {
			return 1
		}
	}
	return 0
}

// String base64 representation of Hash
func (hash0 Hash) String() string {
	return base64.StdEncoding.EncodeToString(hash0[:])
}

// Hex representation of Hash
func (hash0 Hash) Hex() string {
	return hex.EncodeToString(hash0[:])
}

// PrintDistancesToHash print the distance from a list of node addresses to
// a hash
func (hash0 Hash) PrintDistancesToHash(list []NodeAddr) {
	str := "Printing distance to " + hash0.String() + ":\n"
	for _, a := range list {
		str += fmt.Sprintln(a.NodeID.XOR(hash0).Hex(), a.Addr)
	}
	fmt.Print(str)
}
