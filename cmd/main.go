// Version 1
/*
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
	// log.Println("requests:", requests)
	// log.Println("responses:", requests)
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
*/

// Version 2
/*

import "allen/jobsearch/companies/kmf/projects/proxy/internal/handler"

func main() {
	// handler
	// service
	// store
	// client

	// Passing requests from the handler to the client via channels can help to decouple the two components and make the code more modular and flexible.

	// By using channels, the handler and client can communicate asynchronously and independently of each other. This can make it easier to handle multiple requests concurrently and can also make the code more fault-tolerant and resilient to errors.

	// Using channels can also make it easier to add new features or modify existing ones. For example, if we wanted to add a caching layer to the client, we could do so without having to modify the handler code. We could simply modify the client code to check the cache before making a request and use a channel to communicate the result back to the handler.

	// Overall, using channels to pass requests from the handler to the client can help to improve the modularity, flexibility, and scalability of the code.

	reqChan := make(chan *handler.Request)
	respChan := make(chan *handler.Response)

	client := client.NewHTTPClient()
	go handler.HandleRequests(client, reqChan, respChan)

	// send requests
	req1 := &handler.Request{
		Method: "GET",
		URL:    "https://example.com",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		},
	}
	req2 := &handler.Request{
		Method: "POST",
		URL:    "https://example.com",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	reqChan <- req1
	reqChan <- req2

	// receive responses
	resp1 := <-respChan
	resp2 := <-respChan

	// handle responses
	// ...
}

// In this example, the handler package defines the Request and Response types, as well as the Client interface and the HandleRequests function. The client package defines the HTTPClient type, which implements the Client interface. The main package creates the channels for passing requests and responses, creates an instance of the HTTPClient, and sends requests to the HandleRequests function via the reqChan channel. The responses are received via the respChan channel and can be handled as needed.

// Algorythm:
// 0. main composes object with
// 1. Handler gets and saves the request and sends it to ch
// 2. Client catches it from ch and makes request, sends response back via ch
// 3. Handler catches the response from ch and makes response

// instead of each handler calling client and back.

// all handlers write to map and send to -> chan
// chan -> client reads from chan (worker pool) runs them
// handlers need to read chan by and select by ID

//						  store
//							^
// handler ->	  			|
// handler -> 		services - > agent
// handler ->

// Direct

// handler <-> agent
// handler <-> agent
// handler <-> agent

// Via channel

// handler -> |---------------|
// handler -> | req, req, req | - agent
// handler -> |_______________|

// TODO
// - explain use channels to

// caching
// Let's say, if the same request comes: -r-r-r-r-
// We can store it in a map [req] res
// We can read from that map using mutex

// throttling
// If 1000 request per second to same server:
//
*/

// main package
package main

import (
	"log"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil { // actually need to check if it's a temporary error
			log.Fatal(err)
		}

		go copyToStderr(conn)

	}
}

// func proxy(conn net.Conn) {

// }

func copyToStderr(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Print(err)
			return
		}
		os.Stderr.Write(buf[:n])
	}
}
