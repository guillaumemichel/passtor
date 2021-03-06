package passtor

import (
	"bytes"
	"errors"
	"sync"
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
	return H(append(SaltToBytes(accountMetaData.SecretSalt), append(NonceToBytes(accountMetaData.PrivateKeySeedNonce), NonceToBytes(accountMetaData.SymmetricKeyNonce)...)...))
}

func HashLogins(logins map[Hash]Login) Hash {
	data := make([]byte, len(logins)*HASHSIZE)

	i := 0
	for _, key := range GetKeysSorted(logins) {
		copy(data[i:], HashToBytes(logins[key].Hash()))
		i += HASHSIZE
	}

	return H(data)
}

func (account Account) GetSignData() []byte {
	return append(HashToBytes(account.ID),
		append(HashToBytes(account.Keys.Hash()),
			append([]byte{byte(account.Version)},
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

func (accountClient AccountClient) GetID() Hash {
	return H([]byte(accountClient.ID))
}

func (accountClient AccountClient) ToEmptyAccount(secret Secret, secretSalt Salt) (Account, error) {
	keys, skNonce, symmKNonce, err := accountClient.Keys.ToKeys(secret)
	if err != nil {
		return Account{}, err
	}

	return Account{
		ID:      accountClient.GetID(),
		Keys:    keys,
		Version: 0,
		Data:    map[Hash]Login{},
		MetaData: AccountMetaData{
			SecretSalt:          secretSalt,
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

func (login Login) ToLoginClient(symmK SymmetricKey) (LoginClient, error) {
	servicePlain, err := Decrypt(login.Service, login.MetaData.ServiceNonce, symmK)
	if err != nil {
		return LoginClient{}, err
	}

	usernamePlain, err := Decrypt(login.Credentials.Username, login.MetaData.UsernameNonce, symmK)
	if err != nil {
		return LoginClient{}, err
	}

	return LoginClient{
		Service:  string(servicePlain),
		Username: string(usernamePlain),
	}, nil
}

func (loginClient LoginClient) GetID(symmK SymmetricKey) Hash {
	return H(append([]byte(loginClient.Service), append([]byte(loginClient.Username), SymmetricKeyToBytes(symmK)...)...))
}

func (loginClient LoginClient) ToNewLogin(keysClient KeysClient, loginPassword string) (Login, error) {
	serviceEncrypted, serviceNonce, err := Encrypt([]byte(loginClient.Service), keysClient.SymmetricKey)
	if err != nil {
		return Login{}, err
	}

	usernameEncrypted, usernameNonce, err := Encrypt([]byte(loginClient.Username), keysClient.SymmetricKey)
	if err != nil {
		return Login{}, err
	}

	password := loginPassword
	if password == "" {
		password, err = Passphrase()
		if err != nil {
			return Login{}, err
		}
	}

	passwordEncrypted, passwordNonce, err := Encrypt([]byte(password), keysClient.SymmetricKey)
	if err != nil {
		return Login{}, err
	}

	return Login{
		ID:      loginClient.GetID(keysClient.SymmetricKey),
		Service: serviceEncrypted,
		Credentials: Credentials{
			Username: usernameEncrypted,
			Password: passwordEncrypted,
		},
		MetaData: LoginMetaData{
			ServiceNonce:  serviceNonce,
			UsernameNonce: usernameNonce,
			PasswordNonce: passwordNonce,
		},
	}, nil
}

func (account Account) AddLogin(loginClient LoginClient, loginPassword string, keysClient KeysClient) (Account, error) {
	if !account.Verify() {
		return Account{}, errors.New("account does not verify")
	}

	login, err := loginClient.ToNewLogin(keysClient, loginPassword)
	if err != nil {
		return Account{}, err
	}

	if _, ok := account.Data[login.ID]; ok {
		return Account{}, errors.New("login already exists")
	}

	newLogins := DuplicateMap(account.Data)
	newLogins[login.ID] = login

	return Account{
		ID:        account.ID,
		Keys:      account.Keys,
		Version:   account.Version + 1,
		Data:      newLogins,
		MetaData:  account.MetaData,
		Signature: Signature{},
	}.Sign(keysClient.PrivateKey), nil
}

func (account Account) DeleteLogin(ID Hash, sk PrivateKey) (Account, error) {
	if !account.Verify() {
		return Account{}, errors.New("account does not verify")
	}

	if _, ok := account.Data[ID]; !ok {
		return Account{}, errors.New("login does not exist")
	}

	newLogins := DuplicateMap(account.Data)
	delete(newLogins, ID)

	return Account{
		ID:        account.ID,
		Keys:      account.Keys,
		Version:   account.Version + 1,
		Data:      newLogins,
		MetaData:  account.MetaData,
		Signature: Signature{},
	}.Sign(sk), nil
}

func (account Account) UpdateLoginPassword(ID Hash, loginPassword string, keysClient KeysClient) (Account, error) {
	if !account.Verify() {
		return Account{}, errors.New("account does not verify")
	}

	if _, ok := account.Data[ID]; !ok {
		return Account{}, errors.New("login does not exists")
	}

	servicePlain, err := Decrypt(account.Data[ID].Service, account.Data[ID].MetaData.ServiceNonce, keysClient.SymmetricKey)
	if err != nil {
		return Account{}, err
	}

	usernamePlain, err := Decrypt(account.Data[ID].Credentials.Username, account.Data[ID].MetaData.UsernameNonce, keysClient.SymmetricKey)
	if err != nil {
		return Account{}, err
	}

	login, err := LoginClient{
		Service:  string(servicePlain),
		Username: string(usernamePlain),
	}.ToNewLogin(keysClient, loginPassword)
	if err != nil {
		return Account{}, err
	}

	newLogins := DuplicateMap(account.Data)
	newLogins[login.ID] = login

	return Account{
		ID:        account.ID,
		Keys:      account.Keys,
		Version:   account.Version + 1,
		Data:      newLogins,
		MetaData:  account.MetaData,
		Signature: Signature{},
	}.Sign(keysClient.PrivateKey), nil
}

func (account Account) GetLoginClientList(symmK SymmetricKey) ([]LoginClient, error) {
	list := make([]LoginClient, len(account.Data))

	i := 0
	for _, id := range GetKeysSorted(account.Data) {
		loginClient, err := account.Data[id].ToLoginClient(symmK)
		if err != nil {
			return nil, err
		}

		list[i] = loginClient
		i++
	}

	return list, nil
}

func (account Account) GetLoginPassword(loginClient LoginClient, symmK SymmetricKey) ([]byte, error) {
	if login, ok := account.Data[loginClient.GetID(symmK)]; ok {
		password, err := Decrypt(login.Credentials.Password, login.MetaData.PasswordNonce, symmK)
		if err != nil {
			return nil, err
		}

		return password, nil
	}

	return nil, errors.New("login does not exist")
}

// Store an account on a node
func (p *Passtor) Store(newAccount Account, repl uint32) error {
	if !newAccount.Verify() {
		return errors.New("account does not verify")
	}

	if oldAccount, ok := p.Accounts[newAccount.ID]; ok {
		oldAccount.Mutex.Lock()
		if newAccount.Version <= oldAccount.Account.Version {
			oldAccount.Mutex.Unlock()
			return errors.New(ALREADYSTORED)
		}
		if bytes.Compare(newAccount.Keys.PublicKey, oldAccount.Account.Keys.PublicKey) != 0 {
			oldAccount.Mutex.Unlock()
			return errors.New("public key changed")
		}

		oldAccount.Account = newAccount
		oldAccount.Repl = repl

		oldAccount.Mutex.Unlock()
	} else {
		p.Accounts[newAccount.ID] = &AccountInfo{
			Account: newAccount,
			Repl:    repl,
			Mutex:   &sync.Mutex{},
		}
	}
	return nil
}

// Delete an account from a node storage
func (p *Passtor) Delete(id Hash) {
	m := p.Accounts[id].Mutex
	//p.Accounts[id].Mutex.Lock()
	m.Lock()
	delete(p.Accounts, id)
	m.Unlock()
}

func (accountNetwork AccountNetwork) ToAccount() Account {
	logins := make(map[Hash]Login, len(accountNetwork.Data))

	for _, login := range accountNetwork.Data {
		logins[login.ID] = login
	}

	return Account{
		ID:        accountNetwork.ID,
		Keys:      accountNetwork.Keys,
		Version:   accountNetwork.Version,
		Data:      logins,
		MetaData:  accountNetwork.MetaData,
		Signature: accountNetwork.Signature,
	}
}

func (account Account) ToAccountNetwork() AccountNetwork {
	logins := make([]Login, len(account.Data))

	i := 0
	for _, login := range account.Data {
		logins[i] = login
		i++
	}

	return AccountNetwork{
		ID:        account.ID,
		Keys:      account.Keys,
		Version:   account.Version,
		Data:      logins,
		MetaData:  account.MetaData,
		Signature: account.Signature,
	}
}
