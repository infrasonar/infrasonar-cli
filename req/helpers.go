package req

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/infrasonar/infrasonar-cli/cli"
)

func httpGet(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"User-Agent": {fmt.Sprintf("InfraSonarCli/%s", cli.Version)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	if err := errForResponse(resp, true); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func httpAuth(method, url, token string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"User-Agent":    {fmt.Sprintf("InfraSonarCli/%s", cli.Version)},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	if err := errForResponse(resp, true); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func httpJsonMore(method, url, token string, data any, detail bool) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"User-Agent":    {fmt.Sprintf("InfraSonarCli/%s", cli.Version)},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
		"Content-Type":  {"application/json"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	if err := errForResponse(resp, detail); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func httpJson(method, url, token string, data any) ([]byte, error) {
	return httpJsonMore(method, url, token, data, true)
}

func errForResponse(resp *http.Response, detail bool) error {
	if resp.StatusCode/100 == 2 {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(resp.Status)
	}
	bodyString := string(body)
	if detail {
		return fmt.Errorf("%s Response: %s", resp.Status, bodyString)
	}
	return errors.New(bodyString)
}
