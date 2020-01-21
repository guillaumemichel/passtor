package main

import (
	"bytes"
	"testing"

	"gitlab.gnugen.ch/gmichel/passtor/crypto"
)

// go test gitlab.gnugen.ch/gmichel/passtor/tests

func getUser() (string, string) {
	return "issou@epfl.ch", "super-strong-master-password"
}

func TestGenerateReturnsKeysOfValidSize(t *testing.T) {

	_, _, symmK, err := crypto.Generate()

	var emptySymmK = crypto.SymmetricKey{}
	if err != nil {
		t.Fail()
	}
	if symmK == emptySymmK {
		t.Fail()
	}
}

func TestDecryptionOfAnEncryptedMessageReturnsTheMessage(t *testing.T) {
	_, _, symmK, _ := crypto.Generate()

	msg, _ := crypto.RandomBytes(1024)

	ciphertext, nonce, _ := crypto.Encrypt(msg, symmK)
	plaintext, _ := crypto.Decrypt(ciphertext, nonce, symmK)

	if bytes.Compare(msg, plaintext) != 0 {
		t.Fail()
	}
}

func TestVerifyOfASignatureReturnsCorrect(t *testing.T) {
	pk, sk, _, _ := crypto.Generate()

	msg, _ := crypto.RandomBytes(1024)

	sig := crypto.Sign(msg, sk)
	correct := crypto.Verify(msg, sig, pk)

	if !correct {
		t.Fail()
	}
}
