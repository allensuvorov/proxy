// package main

// func Test_handleRequest(t *testing.T) {
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handleRequest(tt.args.w, tt.args.r)
// 		})
// 	}
// }

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleRequest(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		method   string
		url      string
		headers  map[string]string
		payload  interface{}
		expected int
	}{
		{
			name:     "Valid request",
			method:   http.MethodGet, // "GET",
			url:      "https://jsonplaceholder.typicode.com/todos/1",
			headers:  map[string]string{"Authorization": "Bearer abc123"},
			payload:  nil,
			expected: http.StatusOK,
		},
		{
			name:     "Invalid method",
			method:   "INVALID",
			url:      "https://jsonplaceholder.typicode.com/todos/1",
			headers:  map[string]string{"Authorization": "Bearer abc123"},
			payload:  nil,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid URL",
			method:   "GET",
			url:      "invalid-url",
			headers:  map[string]string{"Authorization": "Bearer abc123"},
			payload:  nil,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid headers",
			method:   "GET",
			url:      "https://jsonplaceholder.typicode.com/todos/1",
			headers:  map[string]string{"Invalid-Header": "Invalid-Value"},
			payload:  nil,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid payload",
			method:   "GET",
			url:      "https://jsonplaceholder.typicode.com/todos/1",
			headers:  map[string]string{"Authorization": "Bearer abc123"},
			payload:  "invalid-payload",
			expected: http.StatusBadRequest,
		},
	}

	// Create a new test server with the proxy handler
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()

	// Loop through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert payload to JSON
			var jsonPayload []byte
			if tc.payload != nil {
				var err error
				jsonPayload, err = json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("Failed to marshal payload: %v", err)
				}
			}

			// Create a new request with the test case data
			req, err := http.NewRequest(tc.method, ts.URL, bytes.NewBuffer(jsonPayload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			// Send the request and get the response
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Check the response status code
			if resp.StatusCode != tc.expected {
				t.Errorf("Expected status code %d, but got %d", tc.expected, resp.StatusCode)
			}
		})
	}
}
