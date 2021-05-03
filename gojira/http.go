package gojira

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func SendHttpRequest(
	requestMethod string,
	requestUrl string,
	requestBody io.Reader,
	headers map[string]string,
	successfulStatusCode int) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(requestMethod, requestUrl, requestBody)
	if err != nil {
		log.Fatal(err)
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != successfulStatusCode {
		log.Fatalf("There was an error when performing request:\n%s %s\nResponse code was: %d\nResponse body:\n%s", requestMethod, requestUrl, resp.StatusCode, string(body))
	}
	return body
}
