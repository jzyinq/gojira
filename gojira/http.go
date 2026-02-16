package gojira

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func SendHttpRequest(
	requestMethod string,
	requestUrl string,
	requestBody io.Reader,
	headers map[string]string,
	successfulStatusCode int) ([]byte, error) {
	client := &http.Client{}
	logrus.Debugf("sending %s request to %s", requestMethod, requestUrl)
	if requestBody != nil {
		logrus.Debugf("request body: %s", requestBody)
	}
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != successfulStatusCode {
		logrus.Errorf("There was an error when performing request:\n%s %s\nResponse code was: %d\n"+
			"Response body:\n%s", requestMethod, requestUrl, resp.StatusCode, string(body))
		return nil, fmt.Errorf("There was an error when performing request:\n%s %s\nResponse code was: %d\n"+
			"Response body:\n%s", requestMethod, requestUrl, resp.StatusCode, string(body))
	}
	return body, nil
}

// CreateBasicAuthHeaders creates HTTP headers with Basic authentication
func CreateBasicAuthHeaders(username, password string) map[string]string {
	authToken := fmt.Sprintf("%s:%s", username, password)
	authHeader := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(authToken)))
	return map[string]string{
		"Authorization": authHeader,
		"Content-Type":  "application/json",
	}
}

// CreateBearerAuthHeaders creates HTTP headers with Bearer token authentication
func CreateBearerAuthHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
}

// SendJSONRequest marshals the payload to JSON and sends an HTTP request
// Returns the response body or an error
func SendJSONRequest(
	method string,
	url string,
	payload interface{},
	headers map[string]string,
	expectedStatus int) ([]byte, error) {

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	requestBody := bytes.NewBuffer(payloadJSON)
	return SendHttpRequest(method, url, requestBody, headers, expectedStatus)
}
