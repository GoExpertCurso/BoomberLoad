package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/GoExpertCurso/BoomerLoad/internal/entity"
)

var (
	url           = flag.String("url", "", "")
	requests      = flag.Int("requests", 50, "")
	numberThreads = flag.Int("concurrency", 2, "")
)

var usage = `Usage: boomerLoad -url=http://example.com [options]

	Options:
	-r Number of requests to run. Default is 100.
	-c Number of requests to run concurrently. Total number of requests cannot be smaller than the concurrency level. Default is 1.
`

func main() {
	flag.Usage = func() {
		fmt.Printf(usage, runtime.NumCPU())
	}

	flag.Parse()

	if flag.NFlag() < 1 {
		usageAndExit("")
	}

	url := *url
	num := *requests
	conc := *numberThreads

	reg := regexp.MustCompile(`\bhttps?://\S+\b`)

	if url == "" {
		usageAndExit("url is required. ")
	}

	if !reg.MatchString(url) {
		usageAndExit("url is invalid. ")
	}

	if num <= 0 && conc <= 0 {
		usageAndExit("-r and -c cannot be smaller than 1. ")
	}

	if conc <= 0 {
		usageAndExit("-c cannot be smaller than 1. ")
	}
	if num < conc {
		usageAndExit("-r cannot be less than -c.")
	}

	work := entity.NewWorker(url, num, conc)

	work.Done = make(chan bool)
	work.ResultChan = make(chan bool, num)
	start := time.Now()
	fmt.Printf("Starting load test for %s with %d concurrent requests for %d requests\n", url, conc, num)
	for i := 0; i < conc; i++ {
		go work.Worker()
	}
	lista := []entity.RequestDetails{}

	for j := 0; j < num; j++ {
		<-work.ResultChan
		work.CompletedRequests++
		lista = append(lista, *<-work.HttpDetails)
	}

	counter := entity.NewResponseCounter()

	for _, v := range lista {
		counter.IncrementStatusCodeCount(v.Code)
	}

	work.Close()

	durationSeconds := time.Since(start).Seconds()
	requestsPerSecond := float64(work.CompletedRequests) / durationSeconds

	fmt.Printf("Load test completed in %f seconds\n", durationSeconds)
	fmt.Printf("Total requests: %d\n", work.CompletedRequests)
	fmt.Printf("Requests per second: %2.4f\n", requestsPerSecond)
	counter.PrintStatusCodes()
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
