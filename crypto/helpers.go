package crypto

import "crypto/rand"

// RandomBytes generates an array of random bytes of the given size
func RandomBytes(size uint32) ([]byte, error) {

	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// BytesToSymmetricKey creates a symmetric key from an array of bytes
func BytesToSymmetricKey(array []byte) SymmetricKey {

	if len(array) < SYMMKEYLENGTH {
		panic("Array is expected to have size at least 64")
	}

	var fixedSize [SYMMKEYLENGTH]byte = [SYMMKEYLENGTH]byte{}
	for i, b := range array {
		fixedSize[i] = b
	}

	return SymmetricKey(fixedSize)
}

// SymmetricKeyToBytes converts a symmetric key to a raw array of bytes
func SymmetricKeyToBytes(symmK SymmetricKey) []byte {

	array := make([]byte, len(symmK))

	for _, b := range symmK {
		array = append(array, b)
	}

	return array
}
