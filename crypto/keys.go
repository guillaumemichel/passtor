package crypto

import (
	"golang.org/x/crypto/ed25519"
)

// SYMMKEYLENGTH is the length of the symmetric encryption/decryption key
const SYMMKEYLENGTH = 64

// SymmetricKey is a special type to designate symmetric encryption/decryption keys
type SymmetricKey [SYMMKEYLENGTH]byte

// Keys stores the keys belonging to a user
type Keys struct {
	SignPublicKey           ed25519.PublicKey  // public key to verifiy requests
	EncryptedSignPrivateKey ed25519.PrivateKey // private key to sign requests (encrypted under user's secret)
	EncryptedSymmEncKey     *SymmetricKey      // symmetric encryption key (encrypted under user's secret)
}

func generateSymmetricKey(size uint32) (*SymmetricKey, error) {

	key, err := RandomBytes(SYMMKEYLENGTH)
	if err != nil {
		return nil, err
	}

	symmK := BytesToSymmetricKey(key)
	return &symmK, nil

}

// Generate generates a triplet of keys needed to the user to encrypt/decrypt, sign and verify data
func Generate(secret UserSecret) (*Keys, error) {

	signPublicK, signPrivK, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	symmK, err := generateSymmetricKey(SYMMKEYLENGTH)
	if err != nil {
		return nil, err
	}

	return &Keys{
		SignPublicKey:           signPublicK,
		EncryptedSignPrivateKey: signPrivK,
		EncryptedSymmEncKey:     symmK,
	}, nil

}
