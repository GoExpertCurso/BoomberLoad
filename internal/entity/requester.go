package entity

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ResponseCounter struct {
	StatusCodeCounts map[int]int
}

func NewResponseCounter() *ResponseCounter {
	return &ResponseCounter{
		StatusCodeCounts: make(map[int]int),
	}
}

func (rc *ResponseCounter) IncrementStatusCodeCount(statusCode int) {
	rc.StatusCodeCounts[statusCode]++
}

func (rc *ResponseCounter) PrintStatusCodes() {
	fmt.Println("Status Code Counts:")
	for code, count := range rc.StatusCodeCounts {
		fmt.Printf("%d: %d\n", code, count)
	}
}

type RequestDetails struct {
	Code int
}

type Work struct {
	//Url of request
	Url string

	//Number of requests to make
	Requests int

	//Number of concurrent requests to make
	NumberConcurrent int

	//Channel to signal completion
	Done chan bool

	//Channel to signal result
	ResultChan chan bool

	//Channel to signal timeout
	Timeout chan bool

	//Channel to requests details
	HttpDetails chan *RequestDetails

	//Total number of requests completed
	CompletedRequests int

	//To ensure channels are closed only once
	once sync.Once
}

func NewRequestDetails(code int) *RequestDetails {
	return &RequestDetails{
		Code: code,
	}
}

func NewWorker(url string, requests, numberConcurrent int) *Work {
	return &Work{
		Url:              url,
		Requests:         requests,
		NumberConcurrent: numberConcurrent,
		Done:             make(chan bool),
		ResultChan:       make(chan bool, requests),
		Timeout:          make(chan bool),
		HttpDetails:      make(chan *RequestDetails, requests),
	}
}

func (w *Work) Worker() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, w.NumberConcurrent)

	for j := 0; j < w.Requests; j++ {
		select {
		case <-w.Done:
			wg.Wait()
			w.once.Do(func() {
				close(w.HttpDetails)
				close(w.ResultChan)
			})
			return
		case sem <- struct{}{}:
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					<-sem
				}()

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.Url, nil)
				if err != nil {
					log.Printf("Error creating request: %v", err)
					return
				}

				client := &http.Client{
					//CheckRedirect: w.redirectHandler,
				}
				res, err := client.Do(req)
				if err != nil {
					log.Printf("Error making request: %v", err)
				}

				if res != nil {
					defer res.Body.Close()

					w.HttpDetails <- &RequestDetails{Code: res.StatusCode}

					for res.StatusCode >= 300 && res.StatusCode <= 399 {
						redirectURL := res.Header.Get("Location")
						if redirectURL == "" {
							log.Printf("Error: Redirect location not found")
							continue
						}
						log.Println("Redirected to:", res.Header.Get("Location"))

						req.URL, err = req.URL.Parse(redirectURL)
						if err != nil {
							log.Printf("Error parsing redirect URL: %v", err)
							continue
						}

						// Make request again to redirected URL
						res, err = client.Do(req)
						if err != nil {
							log.Printf("Error making redirect request: %v", err.Error())
							continue
						}
						w.HttpDetails <- &RequestDetails{Code: res.StatusCode}
						defer res.Body.Close()
					}
				}
				w.ResultChan <- true
				log.Println("Request completed")
			}()
		}
	}
	wg.Wait()
	w.once.Do(func() {
		close(w.HttpDetails)
		close(w.ResultChan)
	})
}

func (w *Work) Close() {
	close(w.Done)
}

/* func (w *Work) redirectHandler(req *http.Request, via []*http.Request) error {
	for _, r := range via {
		fmt.Println("Redirected to:", r.URL)
	}
	return nil
} */
