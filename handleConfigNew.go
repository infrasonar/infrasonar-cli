package main

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

type TConfigNew struct {
	name   string
	token  string
	api    string
	output string
}

func askConfigName() string {
	var name string
	fmt.Print("Name: ")
	if _, err := fmt.Scanln(&name); err == nil && reConfigName.MatchString(name) {
		return name
	}
	fmt.Println("Invalid configuration name, please enter a valid name")
	return askConfigName()
}

func askToken() string {
	fmt.Print("Token: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	token := string(pass[:])
	if reToken.MatchString(token) {
		return token
	}
	fmt.Println("Invalid token, please enter a correct token")
	return askToken()
}

func handleConfigNew(cmd *TConfigNew) {
	if cmd.name == "" {
		cmd.name = askConfigName()
	}
	if cmd.token == "" {
		cmd.token = askToken()
	}
	config, err := configurations.new(cmd.name, cmd.token, cmd.api, cmd.output)
	exitOnErr(err)
	fmt.Printf("Configuration '%s' created\n", config.Name)
	os.Exit(0)
}
