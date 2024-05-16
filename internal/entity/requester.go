package entity

import (
	"fmt"
	"log"
	"net/http"
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
	Time float64
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
	client := &http.Client{
		//CheckRedirect: w.redirectHandler,
	}
	for j := 0; j < w.Requests/w.NumberConcurrent; j++ {
		select {
		case <-w.Done:
			return
		default:
			req, err := http.NewRequest(http.MethodGet, w.Url, nil)
			if err != nil {
				log.Printf("Error creating request: %v", err)
				continue
			}

			res, _ := client.Do(req)
			w.HttpDetails <- NewRequestDetails(res.StatusCode)
			if res.StatusCode >= 300 && res.StatusCode <= 399 {
				log.Println("Redirected to:", res.Header.Get("Location"))
				redirectURL := res.Header.Get("Location")
				if redirectURL == "" {
					log.Printf("Error: Redirect location not found")
					continue
				}
				req.URL, err = req.URL.Parse(redirectURL)
				if err != nil {
					fmt.Printf("Error parsing redirect URL: %v", err)
					continue
				}
				// Make request again to redirected URL
				res, err = client.Do(req)
				if err != nil {
					log.Printf("Error making redirect request: %v", err)
					continue
				}
				w.HttpDetails <- &RequestDetails{Code: res.StatusCode}
				res.Body.Close()
			}
			res.Body.Close()
			w.ResultChan <- true
			log.Println("Request completed")
		}
	}
}

func (w *Work) Close() {
	close(w.Done)
}

/* func redirectHandler(req *http.Request, via []*http.Request) error {
	for _, r := range via {
		fmt.Println("Redirected to:", r.URL)
	}
	return nil
} */
