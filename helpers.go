package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func errForResponse(resp *http.Response) error {
	if resp.StatusCode/100 == 2 {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(resp.Status)
	}
	bodyString := string(body)
	return fmt.Errorf("%s Response: %s", resp.Status, bodyString)
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	if err := errForResponse(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitOutput(out interface{}, output string) {
	switch output {
	case "yaml":
		out, err := yaml.Marshal(&out)
		exitOnErr(err)
		fmt.Println(string(out[:]))
		os.Exit(0)
	case "json":
		out, err := json.Marshal(out)
		exitOnErr(err)
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
