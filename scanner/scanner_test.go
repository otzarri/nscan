package scanner

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestIsValidPort(t *testing.T) {
	tables := []struct {
		port    int
		isValid bool
	}{
		{-21, false},
		{22, true},
		{80, true},
		{443, true},
		{65555, false},
	}

	for _, table := range tables {
		isValid := isValidPort(table.port)

		if isValid != table.isValid {
			t.Errorf("Port validation failed for %v\n", table.port)
		}
	}
}

func TestScanPort(t *testing.T) {
	tables := []struct {
		protocol string
		address  string
		port     uint16
		timeout  uint
		isOpen   bool
	}{
		{"tcp", "127.0.0.1", 1025, 500, true},
		{"tcp", "127.0.0.1", 5000, 1000, false},
		{"tcp", "127.0.0.1", 25000, 1500, true},
		{"tcp", "127.0.0.1", 50000, 1000, false},
		{"tcp", "127.0.0.1", 65535, 500, true},
	}

	for _, table := range tables {
		if table.isOpen {
			address := fmt.Sprintf("%s:%d", table.address, table.port)
			server, err := net.Listen(table.protocol, address)
			if err != nil {
				t.Fatal(err)
				defer server.Close()
			}
		}
		isOpen := ScanPort(table.protocol, table.address, table.port, table.timeout)

		if isOpen != table.isOpen {
			t.Errorf("Error: Port %s://%s:%d is closed and should be closed\n",
				table.protocol, table.address, table.port)
		}
	}
}

func TestGetWoolValueMapList(t *testing.T) {
	tables := []struct {
		protocol        string
		address         string
		ports           string
		timeout         uint
		jobValueMapList []map[string]string
	}{
		{
			"tcp", "example.com", "22", uint(500),
			[]map[string]string{
				{
					"address":  "example.com",
					"port":     "22",
					"protocol": "tcp",
					"timeout":  "500",
				},
			},
		}, {
			"tcp", "example.com", "80,443", uint(1000),
			[]map[string]string{
				{
					"address":  "example.com",
					"port":     "80",
					"protocol": "tcp",
					"timeout":  "1000",
				}, {
					"address":  "example.com",
					"port":     "443",
					"protocol": "tcp",
					"timeout":  "1000",
				},
			},
		},
	}

	for _, table := range tables {
		jobValueMapList := GetWoolValueMapList(
			table.protocol,
			table.address,
			table.ports,
			table.timeout,
		)

		if !reflect.DeepEqual(jobValueMapList, table.jobValueMapList) {
			t.Errorf("Error: Received jobValueMapList (%v) is different from the expected (%v)\n",
				jobValueMapList, table.jobValueMapList)
		}
	}
}

func TestParsePorts(t *testing.T) {
	tables := []struct {
		portString string
		portList   []uint16
	}{
		{"5432", []uint16{uint16(5432)}},
		{"80,443", []uint16{uint16(80), uint16(443)}},
		{"21-25", []uint16{uint16(21), uint16(22), uint16(23), uint16(24), uint16(25)}},
	}

	for _, table := range tables {
		portList := parsePorts(table.portString)

		if !reflect.DeepEqual(portList, table.portList) {
			t.Errorf("Error: Received portList (%v) is different from the expected (%v)\n",
				portList, table.portList)
		}
	}
}
