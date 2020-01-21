package crypto

import (
	"golang.org/x/crypto/chacha20poly1305"
)


const (
	// NONCESIZE size in bytes for a nonce
	NONCESIZE = chacha20poly1305.NonceSizeX
)

// Nonce format for encryption
type Nonce [NONCESIZE]byte

// EncryptedData generic format
type EncryptedData []byte

// Encrypt encrypts the given data under the given key using ChaCha20 stream cipher
func Encrypt(data []byte, key SymmetricKey) (EncryptedData, Nonce, error) {

	cipher, err := chacha20poly1305.NewX(SymmetricKeyToBytes(key))
	if err != nil {
		return nil, Nonce{}, err
	}

	nonce, err := RandomBytes(chacha20poly1305.NonceSizeX)
	if err != nil {
		return nil, Nonce{}, err
	}

	return cipher.Seal(nil, nonce, data, nil), BytesToNonce(nonce), nil
}

// Decrypt decrypts the given ciphertext under the given key
func Decrypt(ciphertext []byte, nonce Nonce, key SymmetricKey) ([]byte, error) {

	if len(ciphertext) <= chacha20poly1305.NonceSizeX {
		panic("Invalid ciphertext")
	}

	cipher, err := chacha20poly1305.NewX(SymmetricKeyToBytes(key))
	if err != nil {
		return nil, err
	}

	return cipher.Open(nil, NonceToBytes(nonce), ciphertext, nil)
}
