package main

import (
	"flag"
	"fmt"
	"os"

	"../../passtor"
)

func handleNewAccount() (*passtor.AccountClient, *passtor.Account) {

	if PromptUser("It seems that you don't have a username yet. Would you like to create a new Passtör account? [y/n]", []string{"y", "n"}) == "n" {
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

func handleUpdates(account *passtor.AccountClient) {
	LaunchCrazyCLI(account)
}

func queryAccount(ID, node string, completion func(*passtor.ServerResponse)) {

	query := &passtor.ClientMessage{
		Pull: passtor.H([]byte(ID)),
	}
	completion(Request(query, node))
}

func main() {
	node := flag.String("node", "127.0.0.1:8080", "IP and port of node to connect to")
	username := flag.String("username", "", "client username")
	flag.Parse()

	if *username == "" {

		accountClient, _ := handleNewAccount()
		handleUpdates(accountClient)

	} else {

		queryAccount(*username, *node, func(response *passtor.ServerResponse) {

			if response.Status == "ok" && response.Data != nil {

				masterPassword := PromptUser("Enter your master password:", nil)
				secret := passtor.GetSecret(*username, masterPassword)
				accountClient, err := response.Data.ToAccountClient(*username, secret)
				AbortOnError(err, "Wrong password")
				handleUpdates(&accountClient)

			}

		})

	}
}
