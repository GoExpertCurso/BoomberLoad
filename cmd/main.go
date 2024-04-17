package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/GoExpertCurso/BoomerLoad/internal/entity"
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

	work := entity.NewWorker(url, num, conc)

	work.Done = make(chan bool)
	work.ResultChan = make(chan bool, num)
	start := time.Now()
	fmt.Printf("\nStarting load test for %s with %d concurrent requests for %d requests\n ", url, conc, num)
	for i := 0; i < conc; i++ {
		go work.Worker()
	}
	lista := []entity.RequestDetails{}
	count200 := 0
	count400 := 0
	count500 := 0
	count429 := 0

	for j := 0; j < num; j++ {
		fmt.Println("Pegando o resultado")
		<-work.ResultChan
		lista = append(lista, *<-work.HttpDetails)
		//<-work.HttpDetails
		//fmt.Println(<-work.HttpDetails)
	}

	for _, v := range lista {
		fmt.Println("code: ", v.Code)
		if v.Code == 200 {
			count200++
		}
		if v.Code == 400 {
			count400++
		}
		if v.Code == 500 {
			count500++
		}
		if v.Code == 429 {
			count429++
		}
	}

	work.Close()
	/* done := make(chan bool)
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


	close(done) */

	durationSeconds := time.Since(start).Seconds()

	requestsPerSecond := float64(*requests) / durationSeconds
	fmt.Printf("Load test completed in %f seconds\n", durationSeconds)
	fmt.Printf("Total requests: %d\n", num)
	fmt.Printf("Requests per second: %2.2f\n", requestsPerSecond)
	fmt.Printf("200 OK: %d\n", count200)
	fmt.Printf("400 OK: %d\n", count400)
	fmt.Printf("500 OK: %d\n", count500)
	fmt.Printf("429 OK: %d\n", count429)
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

/* func worker(url string, reqPerConc int, done <-chan bool, resultChan chan<- bool) {
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
}*/
