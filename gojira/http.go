package gojira

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func SendHttpRequest(
	requestMethod string,
	requestUrl string,
	requestBody io.Reader,
	headers map[string]string,
	successfulStatusCode int) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(requestMethod, requestUrl, requestBody)
	if err != nil {
		return nil, err
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != successfulStatusCode {
		return nil, fmt.Errorf("There was an error when performing request:\n%s %s\nResponse code was: %d\nResponse body:\n%s", requestMethod, requestUrl, resp.StatusCode, string(body))
	}
	return body, nil
}
