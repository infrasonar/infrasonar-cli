package req

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

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

func httpGetAuth(url, token string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	if err := errForResponse(resp); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

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
