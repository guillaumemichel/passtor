package passtor

import (
	"bytes"
	"errors"
)

func (logInMetaData LoginMetaData) Hash() Hash {
	return H(append(NonceToBytes(logInMetaData.ServiceNonce),
		append(NonceToBytes(logInMetaData.UsernameNonce),
			NonceToBytes(logInMetaData.PasswordNonce)...)...))
}

func (credentials Credentials) Hash() Hash {
	return H(append(credentials.Username, credentials.Password...))
}

func (login Login) Hash() Hash {
	return H(append(HashToBytes(login.ID),
		append(login.Service,
			append(HashToBytes(login.Credentials.Hash()), HashToBytes(login.MetaData.Hash())...)...)...))
}

func (keysClient KeysClient) Hash() Hash {
	return H(append(keysClient.PublicKey, append(keysClient.PrivateKey, SymmetricKeyToBytes(keysClient.SymmetricKey)...)...))
}

func (keys Keys) Hash() Hash {
	return H(append(keys.PublicKey, append(keys.PrivateKeySeed, keys.SymmetricKey...)...))
}

func (accountMetaData AccountMetaData) Hash() Hash {
	return H(append(NonceToBytes(accountMetaData.PrivateKeySeedNonce), NonceToBytes(accountMetaData.SymmetricKeyNonce)...))
}

func HashLogins(logins map[Hash]Login) Hash {
	data := make([]byte, len(logins)*HASHSIZE)

	i := 0
	for _, login := range logins {
		copy(data[i:], HashToBytes(login.Hash()))
		i += HASHSIZE
	}

	return H(data)
}

func (account Account) GetSignData() []byte {
	return append(HashToBytes(account.ID),
		append(HashToBytes(account.Keys.Hash()),
			append([]byte{account.Version},
				append(HashToBytes(HashLogins(account.Data)), HashToBytes(account.MetaData.Hash())...)...)...)...)
}

func (account Account) Sign(sk PrivateKey) Account {
	return Account{
		ID:        account.ID,
		Keys:      account.Keys,
		Version:   account.Version,
		Data:      account.Data,
		MetaData:  account.MetaData,
		Signature: Sign(account.GetSignData(), sk),
	}
}

func (account Account) Verify() bool {
	return Verify(account.GetSignData(), account.Signature, account.Keys.PublicKey)
}

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

func (accountClient AccountClient) ToEmptyAccount(secret Secret) (Account, error) {
	keys, skNonce, symmKNonce, err := accountClient.Keys.ToKeys(secret)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:      H([]byte(accountClient.ID)),
		Keys:    keys,
		Version: 0,
		Data:    map[Hash]Login{},
		MetaData: AccountMetaData{
			PrivateKeySeedNonce: skNonce,
			SymmetricKeyNonce:   symmKNonce,
		},
		Signature: Signature{},
	}.Sign(accountClient.Keys.PrivateKey), nil
}

func (account Account) ToAccountClient(ID string, secret Secret) (AccountClient, error) {
	if !account.Verify() {
		return AccountClient{}, errors.New("account does not verify")
	}

	keysClient, err := account.Keys.ToKeysClient(secret, account.MetaData.PrivateKeySeedNonce, account.MetaData.SymmetricKeyNonce)
	if err != nil {
		return AccountClient{}, err
	}

	return AccountClient{
		ID:   ID,
		Keys: keysClient,
	}, nil
}

// TODO: test
func (accounts Accounts) Store(newAccount Account) error {
	if !newAccount.Verify() {
		return errors.New("account does not verify")
	}

	if oldAccount, ok := accounts[newAccount.ID]; ok {
		if newAccount.Version <= oldAccount.Version {
			return errors.New("version is in the past, update local data")
		}
		if bytes.Compare(newAccount.Keys.PublicKey, oldAccount.Keys.PublicKey) != 0 {
			return errors.New("public key changed")
		}
	}

	accounts[newAccount.ID] = newAccount

	return nil
}
