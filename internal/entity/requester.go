package entity

import (
	"fmt"
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

func NewRequestDetails(code int, time float64) *RequestDetails {
	return &RequestDetails{
		Code: code,
		Time: time,
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
	client := http.Client{}
	for j := 0; j < w.Requests/w.NumberConcurrent; j++ {
		select {
		case <-w.Done:
			return
		default:
			req, err := http.NewRequest(http.MethodGet, w.Url, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v", err)
				continue
			}
			res, err := client.Do(req)
			w.HttpDetails <- NewRequestDetails(res.StatusCode, 1.0)
			if err != nil {
				fmt.Printf("Error making request: %v", err)
				continue
			}
			res.Body.Close()
			w.ResultChan <- true
		}
	}
}

func (w *Work) Close() {
	close(w.Done)
}
