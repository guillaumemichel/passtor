package main

import (
	"bytes"
	"testing"

	"gitlab.gnugen.ch/gmichel/passtor"
)

// go test -v -count=1 gitlab.gnugen.ch/gmichel/passtor/tests

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

func TestStoreNewAccount(t *testing.T) {
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

	accounts := passtor.NewPasstor("issou", "127.0.0.1:5000", 1)

	err := accounts.Store(account, 5)

	if err != nil {
		t.Fail()
	}
	if len(accounts.Accounts) != 1 {
		t.Fail()
	}
	if !accounts.Accounts[account.ID].Account.Verify() {
		t.Fail()
	}
}

func TestStoreIncorrectAccountFails(t *testing.T) {
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
	account.Version = 3
	if account.Verify() {
		t.Fail()
	}

	accounts := passtor.NewPasstor("issou", "127.0.0.1:5000", 1)

	err := accounts.Store(account, 5)

	if err == nil {
		t.Fail()
	}
	if len(accounts.Accounts) != 0 {
		t.Fail()
	}
}

func TestStoreUpdateOldAccount(t *testing.T) {
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

	accounts := passtor.NewPasstor("issou", "127.0.0.1:5000", 1)

	err := accounts.Store(account, 5)

	if err != nil {
		t.Fail()
	}
	if len(accounts.Accounts) != 1 {
		t.Fail()
	}
	if !accounts.Accounts[account.ID].Account.Verify() {
		t.Fail()
	}

	err = accounts.Store(account, 5)

	if err == nil {
		t.Fail()
	}
	if len(accounts.Accounts) != 1 {
		t.Fail()
	}
	if !accounts.Accounts[account.ID].Account.Verify() {
		t.Fail()
	}

	pkNew, _, _, _ := passtor.Generate()
	account.Keys.PublicKey = pkNew

	err = accounts.Store(account, 5)

	if err == nil {
		t.Fail()
	}
	if len(accounts.Accounts) != 1 {
		t.Fail()
	}
	if !accounts.Accounts[account.ID].Account.Verify() {
		t.Fail()
	}
}

func TestToNewLoginDataIntegrity(t *testing.T) {
	pk, sk, symmK, _ := passtor.Generate()

	keysClient := passtor.KeysClient{
		PublicKey:    pk,
		PrivateKey:   sk,
		SymmetricKey: symmK,
	}

	loginClient := passtor.LoginClient{
		Service:  "twitter",
		Username: "@trump",
	}

	login, _ := loginClient.ToNewLogin(keysClient)
	loginPlain, _ := login.ToLoginClient(keysClient.SymmetricKey)

	if loginClient.Service != loginPlain.Service {
		t.Fail()
	}
	if loginClient.Username != loginPlain.Username {
		t.Fail()
	}
}

func TestLoginAddUpdateDelete(t *testing.T) {
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
	if len(account.Data) != 0 {
		t.Fail()
	}
	if account.Version != 0 {
		t.Fail()
	}

	loginClientTwitter := passtor.LoginClient{
		Service:  "twitter",
		Username: "@trump",
	}

	account, _ = account.AddLogin(loginClientTwitter, keysClient)
	if len(account.Data) != 1 {
		t.Fail()
	}
	if account.Version != 1 {
		t.Fail()
	}

	for _, v := range account.Data {
		servicePlain, _ := passtor.Decrypt(v.Service, v.MetaData.ServiceNonce, keysClient.SymmetricKey)
		usernamePlain, _ := passtor.Decrypt(v.Credentials.Username, v.MetaData.UsernameNonce, keysClient.SymmetricKey)
		passwordPlain, _ := passtor.Decrypt(v.Credentials.Password, v.MetaData.PasswordNonce, keysClient.SymmetricKey)

		t.Log("SERVICE: " + string(servicePlain))
		t.Log("USERNAME: " + string(usernamePlain))
		t.Log("PASSWORD: " + string(passwordPlain))
	}

	_, err := account.AddLogin(loginClientTwitter, keysClient)
	if err == nil {
		t.Fail()
	}

	account, _ = account.UpdateLoginPassword(loginClientTwitter.GetID(keysClient.SymmetricKey), keysClient)
	if len(account.Data) != 1 {
		t.Fail()
	}
	if account.Version != 2 {
		t.Fail()
	}

	for _, v := range account.Data {
		servicePlain, _ := passtor.Decrypt(v.Service, v.MetaData.ServiceNonce, keysClient.SymmetricKey)
		usernamePlain, _ := passtor.Decrypt(v.Credentials.Username, v.MetaData.UsernameNonce, keysClient.SymmetricKey)
		passwordPlain, _ := passtor.Decrypt(v.Credentials.Password, v.MetaData.PasswordNonce, keysClient.SymmetricKey)

		t.Log("SERVICE: " + string(servicePlain))
		t.Log("USERNAME: " + string(usernamePlain))
		t.Log("PASSWORD: " + string(passwordPlain))
	}

	account, _ = account.DeleteLogin(loginClientTwitter.GetID(keysClient.SymmetricKey), keysClient.PrivateKey)
	if len(account.Data) != 0 {
		t.Fail()
	}
	if account.Version != 3 {
		t.Fail()
	}

	_, err = account.UpdateLoginPassword(loginClientTwitter.GetID(keysClient.SymmetricKey), keysClient)
	if err == nil {
		t.Fail()
	}

	_, err = account.DeleteLogin(loginClientTwitter.GetID(keysClient.SymmetricKey), keysClient.PrivateKey)
	if err == nil {
		t.Fail()
	}
}

func TestGetLoginClientList(t *testing.T) {
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

	loginClientTwitter := passtor.LoginClient{
		Service:  "twitter",
		Username: "@trump_twitter",
	}

	loginClientMastodon := passtor.LoginClient{
		Service:  "mastodon",
		Username: "@trump_mastodon",
	}

	loginClientReddit := passtor.LoginClient{
		Service:  "reddit",
		Username: "@trump_reddit",
	}

	account, _ = account.AddLogin(loginClientTwitter, keysClient)
	account, _ = account.AddLogin(loginClientMastodon, keysClient)
	account, _ = account.AddLogin(loginClientReddit, keysClient)
	if len(account.Data) != 3 {
		t.Fail()
	}
	if account.Version != 3 {
		t.Fail()
	}

	for _, v := range account.Data {
		servicePlain, _ := passtor.Decrypt(v.Service, v.MetaData.ServiceNonce, keysClient.SymmetricKey)
		usernamePlain, _ := passtor.Decrypt(v.Credentials.Username, v.MetaData.UsernameNonce, keysClient.SymmetricKey)
		passwordPlain, _ := passtor.Decrypt(v.Credentials.Password, v.MetaData.PasswordNonce, keysClient.SymmetricKey)

		t.Log("SERVICE: " + string(servicePlain))
		t.Log("USERNAME: " + string(usernamePlain))
		t.Log("PASSWORD: " + string(passwordPlain))
	}

	loginList, _ := account.GetLoginClientList(keysClient.SymmetricKey)
	if len(loginList) != 3 {
		t.Fail()
	}

	for _, l := range loginList {
		t.Log("SERVICE: " + l.Service)
		t.Log("USERNAME: " + l.Username)
	}
}

func TestGetPassword(t *testing.T) {
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

	loginClientTwitter := passtor.LoginClient{
		Service:  "twitter",
		Username: "@trump_twitter",
	}

	loginClientMastodon := passtor.LoginClient{
		Service:  "mastodon",
		Username: "@trump_mastodon",
	}

	loginClientReddit := passtor.LoginClient{
		Service:  "reddit",
		Username: "@trump_reddit",
	}

	account, _ = account.AddLogin(loginClientTwitter, keysClient)
	account, _ = account.AddLogin(loginClientMastodon, keysClient)
	account, _ = account.AddLogin(loginClientReddit, keysClient)

	loginList, _ := account.GetLoginClientList(keysClient.SymmetricKey)

	for _, loginClient := range loginList {
		password, _ := account.GetLoginPassword(loginClient, keysClient.SymmetricKey)

		t.Log(loginClient)
		t.Log(string(password))
	}
}

func TestAccountNetworkIntegrity(t *testing.T) {
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

	loginClientTwitter := passtor.LoginClient{
		Service:  "twitter",
		Username: "@trump_twitter",
	}

	loginClientMastodon := passtor.LoginClient{
		Service:  "mastodon",
		Username: "@trump_mastodon",
	}

	loginClientReddit := passtor.LoginClient{
		Service:  "reddit",
		Username: "@trump_reddit",
	}

	account, _ = account.AddLogin(loginClientTwitter, keysClient)
	account, _ = account.AddLogin(loginClientMastodon, keysClient)
	account, _ = account.AddLogin(loginClientReddit, keysClient)

	accountNetwork := account.ToAccountNetwork()
	accountBack := accountNetwork.ToAccount()

	if !accountBack.Verify() {
		t.Fail()
	}
}

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
	mostR, thresh := passtor.MostRepresented(accounts, 0)
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
	mostR, thresh := passtor.MostRepresented(accounts, 0)
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
	mostR, thresh := passtor.MostRepresented(accounts, 5)
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
	mostR, thresh := passtor.MostRepresented(accounts, 3)
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
	mostR, thresh := passtor.MostRepresented(accounts, 5)
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
	mostR, thresh := passtor.MostRepresented(accounts, 3)
	if mostR == nil || mostR.Signature != account2.Signature {
		t.Fail()
	}
	if !thresh {
		t.Fail()
	}
}
