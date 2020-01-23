package main

import "../../passtor"

func LaunchCrazyCLI(account *passtor.AccountClient) {

	for {
		action := PromptUser("What would you like to do?\n - Add a new login\n - Update an existing login\n - Delete an existing login")
	}
}
