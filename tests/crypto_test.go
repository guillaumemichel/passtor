package main

import (
	"bytes"
	"testing"

	"../crypto"
)

func getUser() (string, string) {
	return "issou@epfl.ch", "super-strong-master-password"
}

func TestGenerateReturnsKeysOfValidSize(t *testing.T) {

	userID, masterPassword := getUser()
	secret, _ := crypto.Secret(userID, masterPassword)
	keys, _ := crypto.Generate(secret)

	if keys.SignPublicKey == nil {
		t.Fail()
	}

	if keys.EncryptedSignPrivateKey == nil {
		t.Fail()
	}

	symmK, _ := crypto.Decrypt(keys.EncryptedSymmEncKey, crypto.BytesToSymmetricKey(secret))
	if keys.EncryptedSymmEncKey == nil || len(symmK) != crypto.SYMMKEYLENGTH {
		t.Fail()
	}
}

func TestDecryptionOfAnEncryptedMessageReturnsTheMessage(t *testing.T) {

	key, _ := crypto.RandomBytes(crypto.SYMMKEYLENGTH)
	msg, _ := crypto.RandomBytes(1024)
	ciphertext, _ := crypto.Encrypt(msg, crypto.BytesToSymmetricKey(key))

	plaintext, _ := crypto.Decrypt(ciphertext, crypto.BytesToSymmetricKey(key))

	if bytes.Compare(msg, plaintext) != 0 {
		t.Fail()
	}
}
