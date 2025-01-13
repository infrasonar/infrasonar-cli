package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"os"

	"github.com/howeyc/gopass"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/re"
	"github.com/infrasonar/infrasonar-cli/req"
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

func IsArray(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Array
}

func ExitOutput(out any, output string) {
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
		switch x := out.(type) {
		case []string:
			for _, i := range x {
				fmt.Println(i)
			}
		case []int:
			for _, i := range x {
				fmt.Println(i)
			}
		case []float32:
			for _, i := range x {
				fmt.Println(i)
			}
		case []float64:
			for _, i := range x {
				fmt.Println(i)
			}
		case string, int, bool, float32, float64:
			fmt.Println(out)
		default:
			fmt.Fprintln(os.Stderr, "output 'simple' not possible, try -o yaml or -o json")
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Println(out)
	os.Exit(0)
}

func EnsureContainer(api, token string, containerId int) *cli.Container {
	if containerId == 0 {
		cid, err := req.GetContainerId(api, token)
		if err != nil {
			if strings.Contains(err.Error(), "This route only works with a container token") {
				fmt.Fprintln(os.Stderr, "use the --container (-c) argument to specify a container or switch to a container token")
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}
		containerId = cid
	}

	container, err := req.GetContainer(api, token, containerId)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return container
}

func InSlice(haystack []string, needle string) *string {
	f := strings.ToLower(needle)
	for _, s := range haystack {
		if strings.ToLower(s) == f {
			return &s
		}
	}
	return nil
}

func RemoveFromSlice(haystack *[]string, needle string) int {
	out := []string{}
	found := 0
	for _, s := range *haystack {
		if s == needle {
			found += 1
		} else {
			out = append(out, s)
		}
	}
	*haystack = out
	return found
}

func Itob(i int) bool {
	return i != 0
}
