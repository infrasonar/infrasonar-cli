package util

import (
	"encoding/json"
	"fmt"

	"os"

	"github.com/howeyc/gopass"
	"github.com/infrasonar/infrasonar-cli/re"
	"gopkg.in/yaml.v3"
)

func AskConfigName() string {
	var name string
	fmt.Print("Name: ")
	if _, err := fmt.Scanln(&name); err == nil && re.ConfigName.MatchString(name) {
		return name
	}
	fmt.Println("Invalid configuration name, please enter a valid name")
	return AskConfigName()
}

func AskToken() string {
	fmt.Print("Token: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	token := string(pass[:])
	if re.Token.MatchString(token) {
		return token
	}
	fmt.Println("Invalid token, please enter a correct token")
	return AskToken()
}

func ExitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ExitOutput(out interface{}, output string) {
	switch output {
	case "yaml":
		out, err := yaml.Marshal(&out)
		ExitOnErr(err)
		fmt.Println(string(out[:]))
		os.Exit(0)
	case "json":
		out, err := json.Marshal(out)
		ExitOnErr(err)
		fmt.Println(string(out[:]))
		os.Exit(0)
	case "simple":
		out := fmt.Sprintf("%v", out)
		fmt.Println(out)
		os.Exit(0)
	}
	fmt.Println(out)
	os.Exit(0)
}
