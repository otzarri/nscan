package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"gitlab.com/josebamartos/wool"
)

var (
	tables []struct {
		protocol string
		address  string
		port     uint16
		timeout  uint
		status   string
	}
	resultTable  map[string]string
	jobValueMaps wool.ValueMapList
)

// jobFunc is a function implementation expected by `wool`. When a new job is
// allocated, `wool` sends the job to `jobFunc`, where it is expected to be an
// implementation to process the job and return a result. The result will be
// received by `wool` and if `Job.Active` field is `false`, it will be discarded.
// Otherwise, the result will be directed to the function `resFunc`.
/* func jobFunc(job wool.Job) wool.Result { */
/* address := job.Values["address"] */
/* port, _ := strconv.Atoi(job.Values["port"]) */
/* protocol := job.Values["protocol"] */
/* timeout, _ := strconv.Atoi(job.Values["timeout"]) */
/* result := wool.Result{Job: job, Values: make(map[string]string)} */
/*  */
/* if scanner.ScanPort(protocol, address, uint16(port), uint(timeout)) { */
/*     result.Values["status"] = "open" */
/* } else { */
/*     result.Values["status"] = "closed" */
/*     result.Job.Active = false */
/* } */
/* return result */
/* } */

// resFunc is a function implementation expected by `wool`. When a new result is
// returned from `jobFunc`, `wool` sends the result to `resFunc`, where it is
// expected to be an implementation to process the result and return a it. The
// result will be received by `wool` and if `Result.Job.Active` field is
// `false`, it will be discarded. Otherwise, the result will be added to the
// result list returned by the function `wool.Work`, executed into `main()`.
/* func resFunc(result wool.Result) wool.Result { */
/* fmt.Printf("Open port: %v\n", result.Job.Values["port"]) */
/* result.Values["status"] = "processed" */
/* return result */
/* } */

func TestJobFunc(t *testing.T) {
	for i, table := range tables {
		if table.status == "open" {
			address := fmt.Sprintf("%s:%d", table.address, table.port)
			server, err := net.Listen(table.protocol, address)
			if err != nil {
				t.Fatal(err)
				defer server.Close()
			}

			valueMap := wool.ValueMap{
				"protocol": table.protocol,
				"address":  table.address,
				"port":     strconv.Itoa(int(table.port)),
				"timeout":  strconv.Itoa(int(table.timeout)),
				"status":   table.status,
			}

			job := wool.Job{Id: i, Active: true, Values: valueMap}
			result := jobFunc(job)
			server.Close()

			if result.Values["status"] != table.status {
				t.Errorf("Error: Port %s://%s:%d expected to be %s but it is %s\n",
					table.protocol, table.address, table.port, result.Values["status"], table.status)
			}
		}
	}
}

func TestResFunc(t *testing.T) {
	for i, table := range tables {
		if table.status == "open" {
			address := fmt.Sprintf("%s:%d", table.address, table.port)
			server, err := net.Listen(table.protocol, address)
			if err != nil {
				t.Fatal(err)
				defer server.Close()
			}

			valueMap := wool.ValueMap{
				"protocol": table.protocol,
				"address":  table.address,
				"port":     strconv.Itoa(int(table.port)),
				"timeout":  strconv.Itoa(int(table.timeout)),
				"status":   table.status,
			}

			job := wool.Job{Id: i, Active: true, Values: valueMap}
			result := wool.Result{Job: job, Values: map[string]string{}}
			result = resFunc(result)
			server.Close()

			if result.Values["status"] != "processed" {
				t.Errorf("Error: Port %s://%s:%d expected to be %s but it is %s\n",
					table.protocol, table.address, table.port, result.Values["status"], table.status)
			}
		}
	}
}

func TestMainFunc(t *testing.T) {
	args := []string{"-host", "localhost"}
	os.Args = append([]string{"mainArgs"}, args...)
	main()
}

func init() {
	tables = []struct {
		protocol string
		address  string
		port     uint16
		timeout  uint
		status   string
	}{
		{"tcp", "127.0.0.1", 1025, 500, "open"},
		{"tcp", "127.0.0.1", 5000, 1000, "closed"},
		{"tcp", "127.0.0.1", 25000, 1500, "open"},
		{"tcp", "127.0.0.1", 50000, 1000, "closed"},
		{"tcp", "127.0.0.1", 65535, 500, "open"},
	}
}
