// client package
package client

// import (
// 	"allen/jobsearch/companies/kmf/projects/proxy/internal/handler"
// 	"io"
// 	"net/http"
// )

// type HTTPClient struct {
// 	httpClient *http.Client
// }

// func NewHTTPClient() *HTTPClient {
// 	return &HTTPClient{
// 		httpClient: &http.Client{},
// 	}
// }

// func (c *HTTPClient) Do(req *handler.Request) (*handler.Response, error) {
// 	// create http request
// 	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for k, v := range req.Headers {
// 		httpReq.Header.Set(k, v)
// 	}

// 	// send http request
// 	httpResp, err := c.httpClient.Do(httpReq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer httpResp.Body.Close()

// 	// create response
// 	resp := &handler.Response{
// 		ID:      "generated_unique_id",
// 		Status:  httpResp.StatusCode,
// 		Headers: make(map[string]string),
// 	}
// 	for k, v := range httpResp.Header {
// 		resp.Headers[k] = v[0]
// 	}
// 	body, err := io.ReadAll(httpResp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp.Body = body

// 	return resp, nil
// }
