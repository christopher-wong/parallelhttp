package parallelhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	pc := New(&http.Client{}, 10)

	if pc == nil {
		t.Error("failed to instantiate paralllelhttp client")
	}
}

func TestQueueRequest(t *testing.T) {
	pc := New(&http.Client{}, 10)
	results := pc.GetResultsChan()

	reqCount := 10

	for i := 1; i < reqCount; i++ {
		endpoint := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i)

		req, _ := http.NewRequest("GET", endpoint, nil)
		pc.QueueRequest(req)
	}

	for r := 1; r < reqCount; r++ {
		result := <-results

		if result.Err != nil {
			fmt.Println(result.Err.Error())
			continue
		}

		var decoded map[string]interface{}
		json.NewDecoder(result.Response.Body).Decode(&decoded)

		fmt.Println(decoded)

		result.Response.Body.Close()
	}
}
