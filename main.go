package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/otzarri/nscan/scanner"
	"github.com/otzarri/wool"
)

var (
	jobValueMapList *[]map[string]string
	workers         *int
)

// jobFunc is a function implementation expected by `wool`. When a new job is
// allocated, `wool` sends the job to `jobFunc`, where it is expected to be an
// implementation to process the job and return a result. The result will be
// received by `wool` and if `Job.Active` field is `false`, it will be discarded.
// Otherwise, the result will be directed to the function `resFunc`.
func jobFunc(job wool.Job) wool.Result {
	address := job.Values["address"]
	port, _ := strconv.Atoi(job.Values["port"])
	protocol := job.Values["protocol"]
	timeout, _ := strconv.Atoi(job.Values["timeout"])
	result := wool.Result{Job: job, Values: make(map[string]string)}

	if scanner.ScanPort(protocol, address, uint16(port), uint(timeout)) {
		result.Values["status"] = "open"
	} else {
		result.Values["status"] = "closed"
		result.Job.Active = false
	}
	return result
}

// resFunc is a function implementation expected by `wool`. When a new result is
// returned from `jobFunc`, `wool` sends the result to `resFunc`, where it is
// expected to be an implementation to process the result and return a it. The
// result will be received by `wool` and if `Result.Job.Active` field is
// `false`, it will be discarded. Otherwise, the result will be added to the
// result list returned by the function `wool.Work`, executed into `main()`.
func resFunc(result wool.Result) wool.Result {
	fmt.Printf("Open port: %v\n", result.Job.Values["port"])
	result.Values["status"] = "processed"
	return result
}

func main() {
	startTime := time.Now()

	// CLI options and argument parsing
	workers = flag.Int("workers", 100, "Concurrent connection number")
	address := flag.String("host", "", "Target hostname or IP (mandatory)")
	ports := flag.String("ports", "1-1024", "-ports 21: Scan port 21\n"+
		"-ports 21,22,23: Scan ports 21, 22 and 23\n"+
		"-ports 21-433: Scan all ports from 21 to 433\n")
	protocol := flag.String("protocol", "tcp", "Only \"tcp\" supported")
	timeout := flag.Uint("timeout", 1000, "Timeout in ms")
	flag.Parse()

	// Mandatory field validation. Print help is missing.
	if *address == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create the list of job data for wool
	valueMapList := scanner.GetWoolValueMapList(*protocol, *address, *ports, *timeout)
	jobValueMapList = &valueMapList

	// Launch the worker pool
	wool.Work(*workers, *jobValueMapList, jobFunc, resFunc)
	endTime := time.Now()
	chrono := endTime.Sub(startTime)
	fmt.Printf("Scan duration: %f seconds\n", chrono.Seconds())
}
