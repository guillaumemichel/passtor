package crypto

import (
	"golang.org/x/crypto/ed25519"
)

// Sign computes the signature of the given message under the given private key
func Sign(data []byte, key ed25519.PrivateKey) ([]byte, error) {
	return key.Sign(nil, data, nil)
}

// Verify checks that the signature and data are valid under the given public key
func Verify(data []byte, signature []byte, key ed25519.PublicKey) bool {
	return ed25519.Verify(key, data, signature)
}
