package crypto

import (
	"crypto/sha256"
	"fmt"
)

// HASHSIZE size of a hash in byte
const HASHSIZE = sha256.Size

// Hash is a flexible type to handle hashes
type Hash [HASHSIZE]byte

// H hashes the given string value
func H(str string) Hash {
	return sha256.Sum256([]byte(str))
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
	if len(hash0) != len(hash1) {
		fmt.Println("Cannot compare hashes of different sizes")
		return 0
	}
	for i := 0; i < len(hash0); i++ {
		if hash0[i] < hash1[i] {
			return -1
		} else if hash0[i] > hash1[i] {
			return 1
		}
	}
	return 0
}
