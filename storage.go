package passtor

func (keysClient KeysClient) ToKeys(secret Secret) (Keys, Nonce, Nonce, error) {
	skSeedEncrypted, skNonce, err := Encrypt(keysClient.PrivateKey.Seed(), secret)
	if err != nil {
		return Keys{}, Nonce{}, Nonce{}, err
	}

	symmKEncrypted, symmKNonce, err := Encrypt(SymmetricKeyToBytes(keysClient.SymmetricKey), secret)
	if err != nil {
		return Keys{}, Nonce{}, Nonce{}, err
	}

	return Keys{
		PublicKey:      keysClient.PublicKey,
		PrivateKeySeed: skSeedEncrypted,
		SymmetricKey:   symmKEncrypted,
	}, skNonce, symmKNonce, nil
}

func (keys Keys) ToKeysClient(secret Secret,
	privateKeySeedNonce Nonce,
	symmetricKeyNonce Nonce) (KeysClient, error) {

	skSeed, err := Decrypt(keys.PrivateKeySeed, privateKeySeedNonce, secret)
	if err != nil {
		return KeysClient{}, err
	}

	symmKey, err := Decrypt(keys.SymmetricKey, symmetricKeyNonce, secret)
	if err != nil {
		return KeysClient{}, err
	}

	return KeysClient{
		PublicKey:    keys.PublicKey,
		PrivateKey:   SeedToPrivateKey(skSeed),
		SymmetricKey: BytesToSymmetricKey(symmKey),
	}, nil
}

func (accountClient AccountClient) ToEmptyAccount(secret Secret) Account {
	return Account{} // TODO
}

func (account Account) ToAccountClient(ID string, secret Secret) AccountClient {
	return AccountClient{} // TODO
}
