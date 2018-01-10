package proxyhost

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var hosts map[string]string

func init() {
	var err error
	hosts, err = readHosts("./hosts")
	if err != nil {
		fmt.Println("hosts file not found.")
	}
}

//FindIP get ip of host name
func FindIP(host string) (string, bool) {
	val, ok := hosts[host]
	return val, ok
}

func readHosts(hostPath string) (map[string]string, error) {
	hosts := make(map[string]string)
	file, err := os.Open(hostPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 0 {
			continue
		}
		if i := strings.IndexAny("#", line); i >= 0 {
			continue
		}
		f := strings.Fields(line)
		hosts[f[1]] = f[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}
