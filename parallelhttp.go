package parallelhttp

import (
	"fmt"
	"net/http"
)

// Client represents an instance of parallelhttp with functions to interact
// with the package
type Client struct {
	queue      chan *http.Request
	results    chan *Response
	httpClient *http.Client
}

// Response represents an object containing the result of executing a queued
// request
type Response struct {
	ID       string
	Response *http.Response
	Err      error
}

func worker(wid int, httpClient *http.Client, jobs <-chan *http.Request, results chan<- *Response) {
	fmt.Printf("worker %d created\n", wid)
	for j := range jobs {

		res, err := httpClient.Do(j)

		parallelResponse := &Response{
			ID:       j.URL.String(),
			Response: nil,
			Err:      nil,
		}

		if err != nil {
			parallelResponse.Err = err
			results <- parallelResponse
		}

		if res.StatusCode < 200 || res.StatusCode > 399 {
			parallelResponse.Err = fmt.Errorf("bad status code: %d", res.StatusCode)
			results <- parallelResponse
		}

		parallelResponse.Response = res
		results <- parallelResponse
	}
}

// New returns a new parallelhttp client
func New(httpClient *http.Client, workerCount int) *Client {
	q := make(chan *http.Request)
	results := make(chan *Response)

	for i := 1; i < workerCount; i++ {
		go worker(i, httpClient, q, results)
	}

	return &Client{
		queue:      q,
		results:    results,
		httpClient: httpClient,
	}
}

// QueueRequest takes an &http.Request{} and adds it to the work list to be
// executed
func (c *Client) QueueRequest(req *http.Request) {
	c.queue <- req
}

// GetResultsChan returns the channel with results
func (c *Client) GetResultsChan() <-chan *Response {
	return c.results
}
