# parallelhttp

This is a basic package that helps abstract away the details of parallelizing
http requests using channels and goroutines.

## Example usage

```golang
import github.com/christopher-wong/parallelhttp
```

This example is also implemented in `TestQueueRequest`

```golang
// create a new parallelhttp client and set the max number of workers
pc := New(&http.Client{}, 10)
// grab the results channel that the executed responses will be returned on
results := pc.GetResultsChan()

// queue some requests
reqCount := 10
for i := 1; i < reqCount; i++ {
    endpoint := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i)

    req, _ := http.NewRequest("GET", endpoint, nil)
    pc.QueueRequest(req)
}

// get the requests back from the results channel
for r := 1; r < reqCount; r++ {
    result := <-results

    fmt.Printf("retrieved \t %s \t %d\n", result.ID, result.Response.StatusCode)
}
```

## TODO

- [ ] implement rate limiting on top of workers
- [ ] find some way of abstracting the channels even further
- [ ] allow usage of alternative http clients like [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp) (generics?)