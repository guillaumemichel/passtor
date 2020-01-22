package main

import (
	"bytes"
	"gitlab.gnugen.ch/gmichel/passtor"
	"gitlab.gnugen.ch/gmichel/passtor/crypto"
	"testing"
)

func TestToKeysToKeysClient(t *testing.T) {
	username, mpass := getUser()

	secret := crypto.GetSecret(username, mpass)
	pk, sk, symmK, _ := crypto.Generate()

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
