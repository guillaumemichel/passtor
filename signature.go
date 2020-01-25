package passtor

import (
	"golang.org/x/crypto/ed25519"
)

const (
	// SIGNATURESIZE size in bytes for a signature
	SIGNATURESIZE = ed25519.SignatureSize
)

// Signature format
type Signature [SIGNATURESIZE]byte

// Sign computes the signature of the given message under the given private key
func Sign(data []byte, key PrivateKey) Signature {
	return BytesToSignature(ed25519.Sign(key, data))
}

// Verify checks that the signature and data are valid under the given public key
func Verify(data []byte, signature Signature, key PublicKey) bool {
	return ed25519.Verify(key, data, SignatureToBytes(signature))
}
