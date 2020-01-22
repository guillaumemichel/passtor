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

	// SECRETLENGTH is the length of the key produced by Argon2
	SECRETLENGTH = SYMMKEYSIZE
)

// Secret is the secret used by the user to locally decrypt its symmetric key K and secret key sk
type Secret = SymmetricKey

// Secret generates a secret for the given user
func GetSecret(userID, masterPassword string) Secret {
	return KDFToSecret(argon2.Key(
		[]byte(userID+masterPassword),
		generateSalt(SALTLENGTH),
		ARGONITERATIONS, ARGONMEMORY,
		ARGONPARALELLISM,
		SECRETLENGTH,
	))
}

// generateSalt generates a default salt
func generateSalt(size uint32) []byte {
	return make([]byte, size)
}
