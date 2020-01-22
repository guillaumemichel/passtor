package passtor

import "gitlab.gnugen.ch/gmichel/passtor/crypto"

func (keysClient KeysClient) ToKeys(secret crypto.Secret) (Keys, crypto.Nonce, crypto.Nonce, error) {
	skSeedEncrypted, skNonce, err := crypto.Encrypt(keysClient.PrivateKey.Seed(), secret)
	if err != nil {
		return Keys{}, crypto.Nonce{}, crypto.Nonce{}, err
	}

	symmKEncrypted, symmKNonce, err := crypto.Encrypt(crypto.SymmetricKeyToBytes(keysClient.SymmetricKey), secret)
	if err != nil {
		return Keys{}, crypto.Nonce{}, crypto.Nonce{}, err
	}

	return Keys{
		PublicKey:      keysClient.PublicKey,
		PrivateKeySeed: skSeedEncrypted,
		SymmetricKey:   symmKEncrypted,
	}, skNonce, symmKNonce, nil
}

func (keys Keys) ToKeysClient(secret crypto.Secret,
	privateKeySeedNonce crypto.Nonce,
	symmetricKeyNonce crypto.Nonce) (KeysClient, error) {

	skSeed, err := crypto.Decrypt(keys.PrivateKeySeed, privateKeySeedNonce, secret)
	if err != nil {
		return KeysClient{}, err
	}

	symmKey, err := crypto.Decrypt(keys.SymmetricKey, symmetricKeyNonce, secret)
	if err != nil {
		return KeysClient{}, err
	}

	return KeysClient{
		PublicKey:    keys.PublicKey,
		PrivateKey:   crypto.SeedToPrivateKey(skSeed),
		SymmetricKey: crypto.BytesToSymmetricKey(symmKey),
	}, nil
}

func (accountClient AccountClient) ToEmptyAccount(secret crypto.Secret) Account {
	return Account{} // TODO
}

func (account Account) ToAccountClient(ID string, secret crypto.Secret) AccountClient {
	return AccountClient{} // TODO
}
