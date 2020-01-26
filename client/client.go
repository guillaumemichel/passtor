package main

import (
	"github.com/rivo/tview"
	"gitlab.gnugen.ch/gmichel/passtor"
)

var client = passtor.Client{}

func pushAccount() { // TODO: do after updating account
	accountNetwork := client.Account.ToAccountNetwork()
	message := &passtor.ClientMessage{
		Push: &accountNetwork,
	}

	response := Request(message, client.Node)
	if response.Status != "ok" {
		FailWithError("Error while pushing your account", response.Debug)
	}
}

func createAccount(username string, masterpass string) {
	pk, sk, symmK, err := passtor.Generate()
	AbortOnError(err, "Unable to generate keys")

	accountClient := passtor.AccountClient{
		ID: username,
		Keys: passtor.KeysClient{
			PublicKey:    pk,
			PrivateKey:   sk,
			SymmetricKey: symmK,
		},
	}

	account, err := accountClient.ToEmptyAccount(passtor.GetSecret(username, masterpass))

	client.AccountClient = accountClient
	client.Account = account

	pushAccount()
}

func downloadAccount(username string, masterpass string) {
	h := passtor.H([]byte(username))
	query := &passtor.ClientMessage{
		Pull: &h,
	}

	response := Request(query, client.Node)

	if response.Status == "ok" && response.Data != nil {
		account := (*response.Data).ToAccount()

		accountClient, err := account.ToAccountClient(username, passtor.GetSecret(username, masterpass))
		AbortOnError(err, "Wrong password")

		client.AccountClient = accountClient
		client.Account = account

		return
	}

	FailWithError("Error while fetching your account", response.Debug)
}

func goToLogin(loginClient passtor.LoginClient) {

	password, err := client.Account.GetLoginPassword(loginClient, client.AccountClient.Keys.SymmetricKey)
	if err != nil {
		errMsg := err.Error()
		FailWithError("Password is fucked up...", &errMsg)
	}

	loginDisplay := tview.NewModal().
		SetText(loginClient.Service + "\n\n" + loginClient.Username + "\n\n" + string(password)).
		AddButtons([]string{"back", "change", "delete"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "back" {
				goToLogins()
			} else if buttonLabel == "change" {
				account, err := client.Account.UpdateLoginPassword(loginClient.GetID(client.AccountClient.Keys.SymmetricKey),
					client.AccountClient.Keys)
				if err != nil {
					errMsg := err.Error()
					FailWithError("Update is fucked up...", &errMsg)
				}

				client.Account = account
				pushAccount()

				goToLogin(loginClient)
			} else if buttonLabel == "delete" {
				account, err := client.Account.DeleteLogin(loginClient.GetID(client.AccountClient.Keys.SymmetricKey),
					client.AccountClient.Keys.PrivateKey)
				if err != nil {
					errMsg := err.Error()
					FailWithError("Delete is fucked up...", &errMsg)
				}

				client.Account = account
				pushAccount()

				goToLogins()
			}
		})

	loginDisplay.SetBorder(true).SetTitle(" passtör ")

	client.App.SetRoot(loginDisplay, true).SetFocus(loginDisplay)
}

func goToLogins() {
	loginsScreen := tview.NewList()

	loginsScreen.AddItem("new", "press n", 'n', func() {
		goToLoginCreation()
	})
	loginsScreen.AddItem("quit", "press q", 'q', func() {
		client.App.Stop()
	})

	logins, err := client.Account.GetLoginClientList(client.AccountClient.Keys.SymmetricKey)
	if err != nil {
		errMsg := err.Error()
		FailWithError("Keys are fucked up...", &errMsg)
	}

	for _, loginClient := range logins {
		tmpLoginClient := passtor.LoginClient{
			Service:  string([]byte(loginClient.Service)),
			Username: string([]byte(loginClient.Username)),
		}

		loginsScreen.AddItem(loginClient.Service, loginClient.Username, '-', func() {
			goToLogin(tmpLoginClient)
		})
	}

	loginsScreen.SetBorder(true).SetTitle(" passtör ")

	client.App.SetRoot(loginsScreen, true).SetFocus(loginsScreen)
}

func goToLoginCreation() {
	service := ""
	username := ""

	createLoginScreen := tview.NewForm().
		AddInputField("service", "", 65, nil, func(serviceNew string) {
			service = serviceNew
		}).
		AddInputField("username", "", 65, nil, func(usernameNew string) {
			username = usernameNew
		}).
		AddButton("add", func() {
			loginClient := passtor.LoginClient{
				Service:  service,
				Username: username,
			}

			account, err := client.Account.AddLogin(loginClient, client.AccountClient.Keys)
			if err != nil {
				errMsg := err.Error()
				FailWithError("Keys are fucked up...", &errMsg)
			}

			client.Account = account
			pushAccount()

			goToLogins()
		}).
		AddButton("cancel", func() {
			goToLogins()
		}).
		SetButtonsAlign(tview.AlignCenter)

	createLoginScreen.SetBorder(true).SetTitle(" passtör ")

	client.App.SetRoot(createLoginScreen, true).SetFocus(createLoginScreen)
}

func main() {
	USERNAME_GLOBAL := ""
	MASTERPASS_GLOBAL := ""

	client = passtor.Client{
		App:           tview.NewApplication(),
		Node:          "127.0.0.1:6000",
		AccountClient: passtor.AccountClient{},
		Account:       passtor.Account{},
	}

	loginScreen := tview.NewForm().
		AddInputField("node", client.Node, 65, nil, func(node string) {
			client.Node = node
		}).
		AddInputField("username", "", 65, nil, func(username string) {
			USERNAME_GLOBAL = username
		}).
		AddPasswordField("masterpass", "", 65, '*', func(masterpass string) {
			MASTERPASS_GLOBAL = masterpass
		}).
		AddButton("login", func() {
			downloadAccount(USERNAME_GLOBAL, MASTERPASS_GLOBAL)
			USERNAME_GLOBAL = ""
			MASTERPASS_GLOBAL = ""

			goToLogins()
		}).
		AddButton("create", func() {
			createAccount(USERNAME_GLOBAL, MASTERPASS_GLOBAL)
			USERNAME_GLOBAL = ""
			MASTERPASS_GLOBAL = ""

			goToLogins()
		}).
		AddButton("quit", func() {
			client.App.Stop()
		}).
		SetButtonsAlign(tview.AlignCenter)

	loginScreen.SetBorder(true).SetTitle(" passtör ")

	if err := client.App.SetRoot(loginScreen, true).SetFocus(loginScreen).Run(); err != nil {
		panic(err)
	}
}
