package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type Response struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int64             `json:"length"`
}

var (
	requests  = make(map[string]*Request)
	responses = make(map[string]*Response)
	mutex     = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/", handleRequest)
	serverAddress := ":8080"
	log.Println("Serving on port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO decouple handler and client, put into separate:
// - functions - can test
// - packages - can test and scale separately

// how can we use channels here? as a queue

// TODO two ways to pass requests to agent:
// - direct call
// - queue them via ch - if need to manage stream

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handleRequest - start")
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := generateID()
	reqURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// if body not needed, use HEAD instead of GET
	// if req.Method == http.MethodGet {
	// 	req.Method = http.MethodHead
	// 	log.Println("If body not needed, use HEAD instead of GET!", req.Method)
	// }

	proxyReq, err := http.NewRequest(req.Method, reqURL.String(), nil)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	for key, value := range req.Headers {
		proxyReq.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Bad gateway", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	// TODO: pick log Body or copy to oblivian
	// _, err = io.Copy(io.Discard, resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, "Bad gateway", http.StatusBadGateway)
		return
	}

	log.Println(string(respBody))

	respHeaders := make(map[string]string)
	for key, value := range resp.Header {
		respHeaders[key] = value[0]
	}

	response := &Response{
		ID:      id,
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Length:  resp.ContentLength,
	}

	mutex.Lock()
	requests[id] = &req
	responses[id] = response
	log.Println("requests:", requests)
	log.Println("responses:", requests)
	mutex.Unlock()

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
	log.Println("handleRequest - end")
}

func generateID() string {
	return fmt.Sprintf("%d", len(requests)+1)
}
