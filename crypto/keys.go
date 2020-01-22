package crypto

import (
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/ed25519"
)

const (
	// SYMMKEYSIZE size in bytes for a symmetric key
	SYMMKEYSIZE = chacha20poly1305.KeySize
)

// PublicKey type
type PublicKey = ed25519.PublicKey

// PrivateKey type
type PrivateKey = ed25519.PrivateKey

// SymmetricKey format
type SymmetricKey [SYMMKEYSIZE]byte

func generateSymmetricKey() (SymmetricKey, error) {
	key, err := RandomBytes(SYMMKEYSIZE)
	if err != nil {
		return SymmetricKey{}, err
	}

	symmK := BytesToSymmetricKey(key)
	return symmK, nil

}

func Generate() (PublicKey, PrivateKey, SymmetricKey, error) {
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, SymmetricKey{}, err
	}

	symmKey, err := generateSymmetricKey()
	if err != nil {
		return nil, nil, SymmetricKey{}, err
	}

	return pk, sk, symmKey, nil
}

func SeedToPrivateKey(seed []byte) PrivateKey {
	return ed25519.NewKeyFromSeed(seed)
}
