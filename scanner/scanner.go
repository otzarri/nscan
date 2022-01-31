package scanner

import (
	"log"
	"net"
	"regexp"
	"sort"
	"strconv"
	"time"

	"gitlab.com/josebamartos/nscan/conv"
)

// isValidPort returns true if the given `port int` is a valid port number.
func isValidPort(port int) bool {
	if port > 0 && port <= 65535 {
		return true
	}
	return false
}

// ScanPort returns true if the given port address is listening. It is
// exported because it's usually called from outside the `scanner` package.
// This way it can be integrated into the `jobFunc` implementation for `wool`.
func ScanPort(protocol string, hostname string, port uint16, timeout uint) bool {
	address := hostname + ":" + strconv.Itoa(int(port))
	duration := time.Duration(timeout * 1000000)

	conn, err := net.DialTimeout(protocol, address, duration)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// parsePorts reads the value of the `portString` parameter and returns
// `portList`, an uint16 slice containing all the ports to be scanned.
func parsePorts(portString string) []uint16 {
	var reRange = regexp.MustCompile(`(?m)^([0-9]{1,5})-([0-9]{1,5})$`)
	var reCsv = regexp.MustCompile(`(?m)^(\d+(?:\.\d+)*)|(\d+)`)
	var portList []uint16

	// Extract port range
	if len(portList) == 0 {
		result := reRange.FindAllStringSubmatch(portString, -1)
		if len(result) > 0 {
			startPort := conv.StringToInt(result[0][1])
			endPort := conv.StringToInt(result[0][2])

			if startPort > endPort {
				oldStartPort := startPort
				startPort = endPort
				endPort = oldStartPort
			}

			for port := startPort; port <= endPort; port++ {
				if isValidPort(port) {
					portList = append(portList, uint16(port))
				}
			}
		}
	}

	// Extract single port or comma separated ports
	if len(portList) == 0 {
		result := reCsv.FindAllStringSubmatch(portString, -1)
		if len(result) > 0 {
			foundPorts := []int{}
			for _, m := range result {
				if m[1] != "" || m[2] != "" {
					intPort := conv.StringToInt(m[1] + m[2])
					if isValidPort(intPort) {
						foundPorts = append(foundPorts, intPort)
					}
				}
			}
			sort.Ints(foundPorts)
			portList = conv.SliceIntToUint16(foundPorts)
		}
	}

	// The cause of an empty slice is a bad formatted value
	if len(portList) == 0 {
		log.Fatalln("Bad format in ports option")
	}
	return portList
}

// GetWoolValueMapList is a function to generate valid input data to run the
// port scans concurrently using the worker pool of the `wool` package. For each
// port to be scanned, a map with all the information required to a scan the
// port. All maps are returned in a slice variable called `woolValueMapList`.
func GetWoolValueMapList(
	protocol string,
	address string,
	ports string,
	timeout uint,
) []map[string]string {
	woolValueMapList := make([]map[string]string, 0, 0)

	portList := parsePorts(ports)

	for _, port := range portList {
		scan := map[string]string{
			"address":  address,
			"port":     strconv.Itoa(int(port)),
			"protocol": protocol,
			"timeout":  strconv.Itoa(int(timeout)),
		}
		woolValueMapList = append(woolValueMapList, scan)
	}
	return woolValueMapList
}
