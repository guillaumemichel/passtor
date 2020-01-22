package passtor

import (
	"golang.org/x/crypto/sha3"
)

// HASHSIZE size of a hash in byte
const HASHSIZE = 64

// Hash is a flexible type to handle hashes
type Hash [HASHSIZE]byte

// H hashes the given string value
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
