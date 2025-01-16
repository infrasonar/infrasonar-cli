package util

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/infrasonar/infrasonar-cli/cli"
	"github.com/infrasonar/infrasonar-cli/re"
	"github.com/infrasonar/infrasonar-cli/req"

	"gopkg.in/yaml.v3"
)

type Simple interface {
	Out() interface{}
}

func AskConfigName() string {
	var name string
	fmt.Print("Name: ")
	if _, err := fmt.Scanln(&name); err == nil && re.MetaKey.MatchString(name) {
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

func Plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

func Color(format string, a ...any) {
	fmt.Print(color.HiYellowString(format, a...))
}

func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)

	if err != nil {
		response = ""
	}

	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}

	if slices.Contains(okayResponses, response) {
		return true
	} else if slices.Contains(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

func ExitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ExitErr(format string, a ...any) {
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func ExitOk(format string, a ...any) {
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, a...)
	os.Exit(0)
}

func IsArray(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Array
}

func IsIntegral(val float64) bool {
	return val == float64(int(val))
}

func Short(s string, n int) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.Trim(s, "\n")
	s = strings.Split(s, "\n")[0]
	if len(s) <= n {
		return s
	}
	return fmt.Sprintf("%s...", s[:10])
}

func ExitOutput(out any, output string, outFn string) {
	Log(outFn, "Write output...")
	fp := os.Stdout
	if outFn != "" {
		fo, err := os.OpenFile(outFn, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output file '%s'\n", outFn)
			os.Exit(1)
		}
		fp = fo
		defer fp.Close()
	}
	switch output {
	case "yaml":
		out, err := yaml.Marshal(&out)
		ExitOnErr(err)
		fmt.Fprintln(fp, string(out[:]))
	case "json":
		out, err := json.Marshal(out)
		ExitOnErr(err)
		fmt.Fprintln(fp, string(out[:]))
	case "simple":
		if v, ok := out.(Simple); ok {
			out = v.Out()
		}

		switch x := out.(type) {
		case []string:
			for _, i := range x {
				fmt.Fprintln(fp, i)
			}
		case []int:
			for _, i := range x {
				fmt.Fprintln(fp, i)
			}
		case []float32:
			for _, i := range x {
				fmt.Fprintln(fp, i)
			}
		case []float64:
			for _, i := range x {
				fmt.Fprintln(fp, i)
			}
		case string, int, bool, float32, float64:
			fmt.Fprintln(fp, out)
		default:
			fmt.Fprintln(os.Stderr, "output 'simple' not possible, try -o yaml or -o json")
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown output format '%s'\n", output)
		os.Exit(1)
	}
	Log(outFn, "Done.")
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

func Log(outFn string, a ...any) {
	if outFn != "" {
		fmt.Println(a...)
	}
}

func TextI(i int) string {
	switch i {
	case 0:
		return "zero"
	case 1:
		return "one"
	case 2:
		return "two"
	case 3:
		return "three"
	case 4:
		return "four"
	case 5:
		return "five"
	case 6:
		return "six"
	case 7:
		return "seven"
	case 8:
		return "eight"
	case 9:
		return "nine"
	case 10:
		return "ten"
	}
	return fmt.Sprintf("%d", i)
}

func HumanizeDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		return fmt.Sprintf("%s day%s", TextI(days), Plural(days))
	} else if hours > 0 {
		return fmt.Sprintf("%s hour%s", TextI(hours), Plural(hours))
	} else if minutes > 0 {
		return fmt.Sprintf("%s minute%s", TextI(minutes), Plural(minutes))
	} else if seconds >= 10 {
		return fmt.Sprintf("%s second%s", TextI(seconds), Plural(seconds))
	} else {
		return "a few seconds"
	}
}
