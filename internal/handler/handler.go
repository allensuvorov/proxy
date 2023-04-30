// handler package
package handler

// type Request struct {
// 	Method  string
// 	URL     string
// 	Headers map[string]string
// }

// type Response struct {
// 	ID      string
// 	Status  int
// 	Headers map[string]string
// 	Body    []byte
// }

// type Client interface {
// 	Do(req *Request) (*Response, error)
// }

// func HandleRequests(client Client, reqChan <-chan *Request, respChan chan<- *Response) {
// 	for req := range reqChan {
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			// handle error
// 		}
// 		respChan <- resp
// 	}
// }
