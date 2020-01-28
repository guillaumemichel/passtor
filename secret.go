package passtor

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

type Salt [SALTLENGTH]byte

func ComputeSecret(userID, masterPassword string, salt Salt) Secret {
	return KDFToSecret(argon2.Key(
		[]byte(userID+masterPassword),
		SaltToBytes(salt),
		ARGONITERATIONS,
		ARGONMEMORY,
		ARGONPARALELLISM,
		SECRETLENGTH,
	))
}

// Secret generates a secret for the given user
func NewSecret(userID, masterPassword string) (Secret, Salt, error) {
	salt, err := RandomBytes(SALTLENGTH)
	if err != nil {
		return Secret{}, Salt{}, err
	}

	return KDFToSecret(argon2.Key(
		[]byte(userID+masterPassword),
		salt,
		ARGONITERATIONS,
		ARGONMEMORY,
		ARGONPARALELLISM,
		SECRETLENGTH,
	)), BytesToSalt(salt), nil
}
