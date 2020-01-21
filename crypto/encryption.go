package crypto

import "golang.org/x/crypto/chacha20poly1305"

// Encrypt encrypts the given data under the given key using ChaCha20 stream cipher
func Encrypt(data []byte, key SymmetricKey) ([]byte, error) {

	cipher, err := chacha20poly1305.NewX(SymmetricKeyToBytes(key))
	if err != nil {
		return nil, err
	}

	nonce, err := RandomBytes(chacha20poly1305.NonceSizeX)
	if err != nil {
		return nil, err
	}

	return append(nonce, cipher.Seal(nil, nonce, data, nil)...), nil

}

// Decrypt decrypts the given ciphertext under the given key
func Decrypt(ciphertext []byte, key SymmetricKey) ([]byte, error) {

	if len(ciphertext) <= chacha20poly1305.NonceSizeX {
		panic("Invalid ciphertext")
	}

	nonce, encrypted := ciphertext[:chacha20poly1305.NonceSizeX], ciphertext[chacha20poly1305.NonceSizeX:]

	cipher, err := chacha20poly1305.NewX(SymmetricKeyToBytes(key))
	if err != nil {
		return nil, err
	}

	return cipher.Open(nil, nonce, encrypted, nil)
}
