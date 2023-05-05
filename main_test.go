package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	// Start the mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, `{"foo": "bar"}`, string(body))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"result": "ok"}`)
	}))
	defer mockServer.Close()

	// Start the proxy server
	go main()

	// Wait for the proxy to start listening on :8080
	time.Sleep(100 * time.Millisecond)

	// Make a test request to the proxy
	requestBody := `{"method": "POST", "url": "` + mockServer.URL + `/test", "headers": {"Content-Type": "application/json"}, "body": {"foo": "bar"}}`
	resp, err := http.Post("http://localhost:8080", "application/json", strings.NewReader(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the request from proxy to mock against the original test request
	assert.Equal(t, mockServer.URL+"/test", resp.Request.URL.String())
	assert.Equal(t, http.MethodPost, resp.Request.Method)
	assert.Equal(t, "application/json", resp.Request.Header.Get("Content-Type"))
	body, err := io.ReadAll(resp.Request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"foo": "bar"}`, string(body))

	// Check the response from proxy against the response from mock
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"result": "ok"}`, string(respBody))
}
