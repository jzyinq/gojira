package gojira

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendHttpRequest(t *testing.T) { //nolint:funlen
	t.Run("successful GET request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "test-value", r.Header.Get("X-Test-Header"))
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"success": true}`)
		}))
		defer server.Close()

		headers := map[string]string{"X-Test-Header": "test-value"}
		body, err := SendHttpRequest("GET", server.URL, nil, headers, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, `{"success": true}`, string(body))
	})

	t.Run("successful POST request with body", func(t *testing.T) {
		requestPayload := `{"name": "test"}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			bodyBytes, _ := io.ReadAll(r.Body)
			assert.Equal(t, requestPayload, string(bodyBytes))
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id": 123}`)
		}))
		defer server.Close()

		headers := map[string]string{"Content-Type": "application/json"}
		body, err := SendHttpRequest("POST", server.URL, bytes.NewBufferString(requestPayload), headers, http.StatusCreated)
		assert.NoError(t, err)
		assert.Equal(t, `{"id": 123}`, string(body))
	})

	t.Run("successful PUT request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"updated": true}`)
		}))
		defer server.Close()

		body, err := SendHttpRequest("PUT", server.URL, nil, map[string]string{}, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, `{"updated": true}`, string(body))
	})

	t.Run("successful DELETE request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		body, err := SendHttpRequest("DELETE", server.URL, nil, map[string]string{}, http.StatusNoContent)
		assert.NoError(t, err)
		assert.Empty(t, body)
	})

	t.Run("handles invalid URL", func(t *testing.T) {
		body, err := SendHttpRequest("GET", "://invalid-url", nil, map[string]string{}, http.StatusOK)
		assert.Error(t, err)
		assert.Nil(t, body)
	})

	errorStatusTests := []struct {
		name       string
		statusCode int
		body       string
		errSnippet string
	}{
		{"handles non-success status code", http.StatusBadRequest, `{"error": "bad request"}`, "Response code was: 400"},
		{"handles 404 error", http.StatusNotFound, `{"error": "not found"}`, "Response code was: 404"},
		{"handles 500 error", http.StatusInternalServerError, `{"error": "internal server error"}`, "Response code was: 500"},
	}
	for _, tc := range errorStatusTests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprint(w, tc.body)
			}))
			defer server.Close()

			body, err := SendHttpRequest("GET", server.URL, nil, map[string]string{}, http.StatusOK)
			assert.Error(t, err)
			assert.Nil(t, body)
			assert.Contains(t, err.Error(), tc.errSnippet)
		})
	}
}
