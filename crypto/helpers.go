package crypto

import "crypto/rand"

// RandomBytes generates an array of random bytes of the given size
func RandomBytes(size uint) ([]byte, error) {

	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// BytesToNonce converts a byte array to a Nonce type.
func BytesToNonce(array []byte) Nonce {
	if len(array) != NONCESIZE {
		panic("Array is expected to have size " + string(NONCESIZE))
	}

	var nonce = Nonce{}
	for i, b := range array {
		nonce[i] = b
	}

	return nonce
}

func NonceToBytes(nonce Nonce) []byte {
	array := make([]byte, len(nonce))

	for i, b := range nonce {
		array[i] = b
	}

	return array
}

// BytesToSymmetricKey creates a symmetric key from an array of bytes
func BytesToSymmetricKey(array []byte) SymmetricKey {
	if len(array) != SYMMKEYSIZE {
		panic("Array is expected to have size " + string(SYMMKEYSIZE))
	}

	var symmKey = [SYMMKEYSIZE]byte{}
	for i, b := range array {
		symmKey[i] = b
	}

	return symmKey
}

// SymmetricKeyToBytes converts a symmetric key to a raw array of bytes
func SymmetricKeyToBytes(symmK SymmetricKey) []byte {

	array := make([]byte, len(symmK))

	for i, b := range symmK {
		array[i] = b
	}

	return array
}

func BytesToSignature(array []byte) Signature {
	if len(array) != SIGNATURESIZE {
		panic("Array is expected to have size " + string(SIGNATURESIZE))
	}

	var sig = Signature{}
	for i, b := range array {
		sig[i] = b
	}

	return sig
}

func SignatureToBytes(signature Signature) []byte {
	array := make([]byte, len(signature))

	for i, b := range signature {
		array[i] = b
	}

	return array
}

func KDFToSecret(array []byte) Secret {
	if len(array) != SECRETLENGTH {
		panic("Array is expected to have size " + string(SECRETLENGTH))
	}

	var secret = Secret{}
	for i, b := range array {
		secret[i] = b
	}

	return secret
}
