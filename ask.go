package main

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

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
	} else {
		fmt.Println("Invalid token, please enter a correct token")
		return askToken()
	}
}

func askProjectName() string {
	fmt.Print("Token: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	token := string(pass[:])
	if reToken.MatchString(token) {
		return token
	} else {
		fmt.Println("Invalid token, please enter a correct token")
		return askToken()
	}
}
