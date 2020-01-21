package crypto

import (
	"golang.org/x/crypto/argon2"
)

const (
	// SALTLENGTH is the length of the salt to be used by Argon2 as KDF
	SALTLENGTH = 16

	// ARGONITERATIONS is the number of iterations to be used in the Argon2 algo
	ARGONITERATIONS = 3

	// ARGONMEMORY is the size in bytes to be used by Argon2 in memory
	ARGONMEMORY = 64 * 1024

	// ARGONPARALELLISM is the number of cores to be used by Argon2
	ARGONPARALELLISM = 2

	// ARGONKEYLENGTH is the length of the key produced by Argon2
	ARGONKEYLENGTH = 32
)

// UserSecret is the secret used by the user to locally decrypt its symmetric key K and private key pk
type UserSecret []byte

// Secret generates a secret for the given user
func Secret(userID, masterPassword string) (UserSecret, error) {

	return argon2.Key(
		[]byte(userID+masterPassword),
		generateSalt(SALTLENGTH),
		ARGONITERATIONS, ARGONMEMORY,
		ARGONPARALELLISM,
		ARGONKEYLENGTH,
	), nil
}

// generateSalt generates a default salt
func generateSalt(size uint32) []byte {
	return make([]byte, size)
}
