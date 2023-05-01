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
	reqRes = make(map[string]*Response)
	mutex  = &sync.Mutex{}
)

func main() {
	http.HandleFunc("/", handleRequest)
	serverAddress := ":8080"
	log.Println("Serving on port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// TODO
// DONE - caching in the map
// - can we move caching to a seperate func?
// - but can we move reqRes so it's not a gloval var?
// - decouple handler and client, put into separate functions - can test
// - context
// - timeouts (as of now open connections stay open forever)

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

	// DIY caching
	mutex.Lock()
	cachedResponce, ok := reqRes[string(body)]
	mutex.Unlock()

	if ok {
		log.Println("request already exists")

		jsonResponse, err := json.Marshal(cachedResponce)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(jsonResponse)
		return
	}

	id := generateID()
	reqURL, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// if response body not needed, use HEAD instead of GET
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

	// TODO: pick to log resp.Body or copy to oblivian
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
	reqRes[string(body)] = response
	log.Println("reqRes:", reqRes)
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
	return fmt.Sprintf("%d", len(reqRes)+1)
}
