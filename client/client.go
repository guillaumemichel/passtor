package main

import (
	"flag"
	"fmt"
	"github.com/rivo/tview"
	"gitlab.gnugen.ch/gmichel/passtor"
	"os"
)

func handleNewAccount() (*passtor.AccountClient, *passtor.Account) {

	if PromptUser("It seems that you don't have a username yet. Would you like to create a new Passtör account? [y/n]",
		[]string{"y", "n"}) == "n" {
		fmt.Println("Okay bye.")
		os.Exit(0)
	}

	username := PromptUser("What would you use as username? ", nil)
	masterPassword := PromptUser("What is your master password?", nil)

	secret := passtor.GetSecret(username, masterPassword)
	pk, sk, k, err := passtor.Generate()
	AbortOnError(err, "Unable to generate keys")

	accountClient := passtor.AccountClient{
		ID: username,
		Keys: passtor.KeysClient{
			PublicKey:    pk,
			PrivateKey:   sk,
			SymmetricKey: k,
		},
	}

	account, err := accountClient.ToEmptyAccount(secret)
	AbortOnError(err, "Unable to create account")
	return &accountClient, &account

}

func handleUpdates(accountClient passtor.AccountClient) {
	LaunchCrazyCLI(accountClient)
}

func queryAccount(ID, node string, completion func(*passtor.ServerResponse)) {

	h := passtor.H([]byte(ID))
	query := &passtor.ClientMessage{
		Pull: &h,
	}

	completion(Request(query, node))
}

func pushAccount(account passtor.Account, node string) {
	accountNetwork := account.ToAccountNetwork()
	message := &passtor.ClientMessage{
		Push: &accountNetwork,
	}

	res := Request(message, node)
	if res.Status != "ok" {
		FailWithError("Error while pushing your account", res.Debug)
	}
}

func createAccount(username string, masterpass string, node string) (passtor.AccountClient, passtor.Account) {
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

	pushAccount(account, node)

	return accountClient, account
}

func downloadAccount(username string, masterpass string, node string) (passtor.AccountClient, passtor.Account) {
	h := passtor.H([]byte(username))
	query := &passtor.ClientMessage{
		Pull: &h,
	}

	response := Request(query, node)

	account := (*response.Data).ToAccount()

	if response.Status == "ok" && response.Data != nil {
		accountClient, err := account.ToAccountClient(username, passtor.GetSecret(username, masterpass))
		AbortOnError(err, "Wrong password")

		return accountClient, account
	}

	FailWithError("Error while fetching your account", response.Debug)

	return passtor.AccountClient{}, passtor.Account{}
}

func goToLogin(app *tview.Application, accountClient passtor.AccountClient, account passtor.Account,
	loginClient passtor.LoginClient, node string) {

	password, err := account.GetLoginPassword(loginClient, accountClient.Keys.SymmetricKey)
	if err != nil {
		errMsg := err.Error()
		FailWithError("Password is fucked up...", &errMsg)
	}

	loginDisplay := tview.NewModal().
		SetText(loginClient.Service + "\n\n" + loginClient.Username + "\n\n" + string(password)).
		AddButtons([]string{"back", "change", "delete"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "back" {
				goToLogins(app, accountClient, account, node)
			} else if buttonLabel == "change" {
				account, err := account.UpdateLoginPassword(loginClient.GetID(accountClient.Keys.SymmetricKey), accountClient.Keys)
				if err != nil {
					errMsg := err.Error()
					FailWithError("Update is fucked up...", &errMsg)
				}

				pushAccount(account, node)
				goToLogin(app, accountClient, account, loginClient, node)
			} else if buttonLabel == "delete" {
				account, err := account.DeleteLogin(loginClient.GetID(accountClient.Keys.SymmetricKey), accountClient.Keys.PrivateKey)
				if err != nil {
					errMsg := err.Error()
					FailWithError("Delete is fucked up...", &errMsg)
				}

				pushAccount(account, node)
				goToLogins(app, accountClient, account, node)
			}
		})

	loginDisplay.SetBorder(true).SetTitle(" passtör ")

	app.SetRoot(loginDisplay, true).SetFocus(loginDisplay)
}

func goToLogins(app *tview.Application, accountClient passtor.AccountClient, account passtor.Account, node string) {
	loginsScreen := tview.NewList()

	loginsScreen.AddItem("new", "press n", 'n', func() {
		goToLoginCreation(app, accountClient, account, node)
	})
	loginsScreen.AddItem("quit", "press q", 'q', func() {
		app.Stop()
	})

	logins, err := account.GetLoginClientList(accountClient.Keys.SymmetricKey)
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
			goToLogin(app, accountClient, account, tmpLoginClient, node)
		})
	}

	loginsScreen.SetBorder(true).SetTitle(" passtör ")

	app.SetRoot(loginsScreen, true).SetFocus(loginsScreen)
}

func goToLoginCreation(app *tview.Application, accountClient passtor.AccountClient, account passtor.Account, node string) {
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

			account, err := account.AddLogin(loginClient, accountClient.Keys)
			if err != nil {
				errMsg := err.Error()
				FailWithError("Keys are fucked up...", &errMsg)
			}

			pushAccount(account, node)
			goToLogins(app, accountClient, account, node)
		}).
		AddButton("cancel", func() {
			goToLogins(app, accountClient, account, node)
		}).
		SetButtonsAlign(tview.AlignCenter)

	createLoginScreen.SetBorder(true).SetTitle(" passtör ")

	app.SetRoot(createLoginScreen, true).SetFocus(createLoginScreen)
}

func main() {
	app := tview.NewApplication()

	NODE_GLOBAL := "127.0.0.1:8080"
	USERNAME_GLOBAL := ""
	MASTERPASS_GLOBAL := ""
	ACCOUNTCLIENT_GLOBAL := passtor.AccountClient{}
	ACCOUNT_GLOBAL := passtor.Account{}

	loginScreen := tview.NewForm().
		AddInputField("node", NODE_GLOBAL, 65, nil, func(node string) {
			NODE_GLOBAL = node
		}).
		AddInputField("username", "", 65, nil, func(username string) {
			USERNAME_GLOBAL = username
		}).
		AddPasswordField("masterpass", "", 65, '*', func(masterpass string) {
			MASTERPASS_GLOBAL = masterpass
		}).
		AddButton("login", func() {
			ACCOUNTCLIENT_GLOBAL, ACCOUNT_GLOBAL = downloadAccount(USERNAME_GLOBAL, MASTERPASS_GLOBAL, NODE_GLOBAL)
			MASTERPASS_GLOBAL = ""

			goToLogins(app, ACCOUNTCLIENT_GLOBAL, ACCOUNT_GLOBAL, NODE_GLOBAL)
		}).
		AddButton("create", func() {
			ACCOUNTCLIENT_GLOBAL, ACCOUNT_GLOBAL = createAccount(USERNAME_GLOBAL, MASTERPASS_GLOBAL, NODE_GLOBAL)
			MASTERPASS_GLOBAL = ""

			goToLogins(app, ACCOUNTCLIENT_GLOBAL, ACCOUNT_GLOBAL, NODE_GLOBAL)
		}).
		AddButton("quit", func() {
			app.Stop()
		}).
		SetButtonsAlign(tview.AlignCenter)

	loginScreen.SetBorder(true).SetTitle(" passtör ")

	if err := app.SetRoot(loginScreen, true).SetFocus(loginScreen).Run(); err != nil {
		panic(err)
	}
}

func mainOLD() {
	node := flag.String("node", "127.0.0.1:8080", "IP and port of node to connect to")
	username := flag.String("username", "", "client username")
	flag.Parse()

	if *username == "" {

		accountClient, account := handleNewAccount()
		accountNetwork := account.ToAccountNetwork()
		message := &passtor.ClientMessage{
			Push: &accountNetwork,
		}

		res := Request(message, *node)
		if res.Status == "ok" {

			handleUpdates(*accountClient)

		} else {

			FailWithError("Error while creating your account", res.Debug)

		}

	} else {

		queryAccount(*username, *node, func(response *passtor.ServerResponse) {

			if response.Status == "ok" && response.Data != nil {

				account := (*response.Data).ToAccount()

				masterPassword := PromptUser("Enter your master password:", nil)
				secret := passtor.GetSecret(*username, masterPassword)
				accountClient, err := account.ToAccountClient(*username, secret)
				AbortOnError(err, "Wrong password")
				handleUpdates(accountClient)

			} else {

				FailWithError("Error while fetching your account", response.Debug)

			}
		})

	}
}
