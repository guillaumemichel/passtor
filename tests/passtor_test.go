package main

import (
	"bytes"
	"gitlab.gnugen.ch/gmichel/passtor"
	"testing"
)

// go test gitlab.gnugen.ch/gmichel/passtor/tests

func getUser() (string, string) {
	return "issou@epfl.ch", "super-strong-master-password"
}

func TestGenerateReturnsKeysOfValidSize(t *testing.T) {

	_, _, symmK, err := passtor.Generate()

	var emptySymmK = passtor.SymmetricKey{}
	if err != nil {
		t.Fail()
	}
	if symmK == emptySymmK {
		t.Fail()
	}
}

func TestDecryptionOfAnEncryptedMessageReturnsTheMessage(t *testing.T) {
	_, _, symmK, _ := passtor.Generate()

	msg, _ := passtor.RandomBytes(1024)

	ciphertext, nonce, _ := passtor.Encrypt(msg, symmK)
	plaintext, _ := passtor.Decrypt(ciphertext, nonce, symmK)

	if bytes.Compare(msg, plaintext) != 0 {
		t.Fail()
	}
}

func TestVerifyOfASignatureReturnsCorrect(t *testing.T) {
	pk, sk, _, _ := passtor.Generate()

	msg, _ := passtor.RandomBytes(1024)

	sig := passtor.Sign(msg, sk)
	correct := passtor.Verify(msg, sig, pk)

	if !correct {
		t.Fail()
	}
}

func TestToKeysToKeysClient(t *testing.T) {
	username, mpass := getUser()

	secret := passtor.GetSecret(username, mpass)
	pk, sk, symmK, _ := passtor.Generate()

	keysClient := passtor.KeysClient{
		PublicKey:    pk,
		PrivateKey:   sk,
		SymmetricKey: symmK,
	}
	keys, nonceSk, nonceSymmK, _ := keysClient.ToKeys(secret)
	keysPlain, _ := keys.ToKeysClient(secret, nonceSk, nonceSymmK)

	if bytes.Compare(keysClient.PublicKey, keysPlain.PublicKey) != 0 {
		t.Fail()
	}

	if bytes.Compare(keysClient.PrivateKey, keysPlain.PrivateKey) != 0 {
		t.Fail()
	}

	if keysClient.SymmetricKey != keysPlain.SymmetricKey {
		t.Fail()
	}
}

func TestAccountSignedVerifies(t *testing.T) {
	username, mpass := getUser()

	secret := passtor.GetSecret(username, mpass)
	pk, sk, symmK, _ := passtor.Generate()

	keysClient := passtor.KeysClient{
		PublicKey:    pk,
		PrivateKey:   sk,
		SymmetricKey: symmK,
	}

	accountClient := passtor.AccountClient{
		ID:   username,
		Keys: keysClient,
	}

	account, _ := accountClient.ToEmptyAccount(secret)

	if !account.Verify() {
		t.Fail()
	}
}

func TestAccountConversionsMatch(t *testing.T) {
	username, mpass := getUser()

	secret := passtor.GetSecret(username, mpass)
	pk, sk, symmK, _ := passtor.Generate()

	keysClient := passtor.KeysClient{
		PublicKey:    pk,
		PrivateKey:   sk,
		SymmetricKey: symmK,
	}

	accountClient := passtor.AccountClient{
		ID:   username,
		Keys: keysClient,
	}

	account, _ := accountClient.ToEmptyAccount(secret)

	accountPlain, _ := account.ToAccountClient(username, secret)

	if accountClient.ID != accountPlain.ID {
		t.Fail()
	}

	if bytes.Compare(accountClient.Keys.PublicKey, accountPlain.Keys.PublicKey) != 0 {
		t.Fail()
	}

	if bytes.Compare(accountClient.Keys.PrivateKey, accountPlain.Keys.PrivateKey) != 0 {
		t.Fail()
	}

	if accountClient.Keys.SymmetricKey != accountPlain.Keys.SymmetricKey {
		t.Fail()
	}
}
