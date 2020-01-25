package main

import (
	"testing"

	"../../passtor"
)

func getAccount() passtor.Account {
	ID := ""
	pass := ""
	secret := passtor.GetSecret(ID, pass)
	pk, sk, k, _ := passtor.Generate()
	keys, _, _, _ := passtor.KeysClient{
		PublicKey:    pk,
		PrivateKey:   sk,
		SymmetricKey: k,
	}.ToKeys(secret)

	return passtor.Account{
		ID:   passtor.H([]byte(ID)),
		Keys: keys,
	}.Sign(sk)
}

func TestMostRepresentedCorrectlyHandlesEmptyArrays(t *testing.T) {
	var accounts []passtor.Account
	mostR, thresh := MostRepresented(accounts, 0)
	if mostR != nil {
		t.Fail()
	}
	if thresh {
		t.Fail()
	}
}

func TestMostRepresentedIgnoresUnverifiedAccounts(t *testing.T) {
	sign, _ := passtor.RandomBytes(passtor.SIGNATURESIZE)
	account := getAccount()
	account.Signature = passtor.BytesToSignature(sign)
	accounts := []passtor.Account{account}
	mostR, thresh := MostRepresented(accounts, 0)
	if mostR != nil {
		t.Fail()
	}
	if thresh {
		t.Fail()
	}
}

func TestMostRepresentedReturnsCorrectAccountIfAllAccountsAreIdenticalAndThreshIsNotMet(t *testing.T) {
	account := getAccount()
	accounts := []passtor.Account{account, account, account, account}
	mostR, thresh := MostRepresented(accounts, 5)
	if mostR == nil || mostR.Signature != account.Signature {
		t.Fail()
	}
	if thresh {
		t.Fail()
	}
}

func TestMostRepresentedReturnsCorrectAccountIfAllAccountsAreIdenticalAndThreshIsMet(t *testing.T) {
	account := getAccount()
	accounts := []passtor.Account{account, account, account, account}
	mostR, thresh := MostRepresented(accounts, 3)
	if mostR == nil || mostR.Signature != account.Signature {
		t.Fail()
	}
	if !thresh {
		t.Fail()
	}
}

func TestMostRepresentedReturnsCorrectAccountAndThreshIsNotMet(t *testing.T) {
	account1 := getAccount()
	account2 := getAccount()
	accounts := []passtor.Account{account2, account2, account1, account2, account1}
	mostR, thresh := MostRepresented(accounts, 5)
	if mostR == nil || mostR.Signature != account2.Signature {
		t.Fail()
	}
	if thresh {
		t.Fail()
	}
}

func TestMostRepresentedReturnsCorrectAccountAndThreshIsMet(t *testing.T) {
	account1 := getAccount()
	account2 := getAccount()
	accounts := []passtor.Account{account2, account2, account1, account2, account1}
	mostR, thresh := MostRepresented(accounts, 3)
	if mostR == nil || mostR.Signature != account2.Signature {
		t.Fail()
	}
	if !thresh {
		t.Fail()
	}
}
