package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	url           = flag.String("url", "", "")
	requests      = flag.Int("r", 50, "")
	numberThreads = flag.Int("c", 200, "")
)

var usage = `Usage: boomerLoad [options...] <url>

	Options:
	-r Number of requests to run. Default is 100.
	-c Number of requests to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 1.
`

func main() {
	flag.Usage = func() {
		fmt.Printf(usage, runtime.NumCPU())
	}
	//flag.Usage()
	fmt.Println()

	flag.Parse()
	/* 	fmt.Println("flag.NArg():", flag.NArg())
	   	if flag.NArg() < 1 {
	   		usageAndExit("")
	   	} */

	url := *url
	num := *requests
	conc := *numberThreads

	if url == "" {
		fmt.Fprint(os.Stderr, "url is required")
	}

	if conc <= 0 {
		usageAndExit("-c cannot be smaller than 1.")
	}
	if num <= 0 || conc <= 0 {
		usageAndExit("-r and -c cannot be smaller than 1.")
	}
	if num < conc {
		usageAndExit("-r cannot be less than -c.")
	}

	done := make(chan bool)
	resultChan := make(chan bool, num)

	start := time.Now()

	fmt.Printf("\nStarting load test for %s with %d concurrent requests for %d requests\n ", url, conc, num)

	reqPerConc := num / conc

	for i := 0; i < conc; i++ {
		go worker(url, reqPerConc, done, resultChan)
	}

	for j := 0; j < num; j++ {
		<-resultChan
	}

	/* for i := 0; i < conc; i++ {
		wg.Add(1)
		for j := 0; j < reqPerConc; j++ {
			fmt.Println("i:", i, "j:", j)
			go func() {
				defer wg.Done()
				client := http.Client{}
				for j := 0; j < reqPerConc; j++ {
					select {
					case <-done:
						return
					default:
						req, err := http.NewRequest("GET", url, nil)
						if err != nil {
							fmt.Printf("Error creating request: %v", err)
							continue
						}
						res, err := client.Do(req)
						if err != nil {
							fmt.Printf("Error making request: %v", err)
							continue
						}
						res.Body.Close()
					}
				}
			}()
		}
	} */

	//time.Sleep(time.Second * time.Duration(1))

	close(done)

	durationSeconds := time.Since(start).Seconds()

	requestsPerSecond := float64(*requests) / durationSeconds
	fmt.Printf("Load test completed in %f seconds\n", durationSeconds)
	fmt.Printf("Total requests: %d\n", num)
	fmt.Printf("Requests per second: %2.2f\n", requestsPerSecond)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprint(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func worker(url string, reqPerConc int, done <-chan bool, resultChan chan<- bool) {
	client := http.Client{}
	for j := 0; j < reqPerConc; j++ {
		fmt.Println("Thread: ", j)
		select {
		case <-done:
			return
		default:
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Printf("Error creating request: %v", err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error making request: %v", err)
				continue
			}
			res.Body.Close()
			resultChan <- true
		}
	}
}
