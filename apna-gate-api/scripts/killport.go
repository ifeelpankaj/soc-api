//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func readPort(envFile string) (string, error) {
	f, err := os.Open(envFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	re := regexp.MustCompile(`^(SERVER_PORT|PORT)=(.*)$`)
	var port string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := re.FindStringSubmatch(line); m != nil {
			port = strings.Trim(strings.TrimSpace(m[2]), `"'`)
		}
	}
	if port == "" {
		return "", fmt.Errorf("no PORT or SERVER_PORT found in %s", envFile)
	}
	return port, nil
}

func killWindows(port string) {
	out, err := exec.Command("cmd", "/C", fmt.Sprintf("netstat -ano | findstr :%s", port)).Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		fmt.Printf("Port %s is free.\n", port)
		return
	}
	killed := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		localAddr, state, pid := fields[1], fields[3], fields[4]
		if !strings.HasSuffix(localAddr, ":"+port) || state != "LISTENING" || killed[pid] {
			continue
		}
		killed[pid] = true
		fmt.Printf("Stopping PID %s...\n", pid)
		exec.Command("taskkill", "/PID", pid, "/F").Run()
	}
	if len(killed) == 0 {
		fmt.Printf("Port %s is free.\n", port)
	}
}

func killUnix(port string) {
	out, err := exec.Command("lsof", "-ti", "tcp:"+port).Output()
	if err != nil || len(strings.TrimSpace(string(out))) == 0 {
		fmt.Printf("Port %s is free.\n", port)
		return
	}
	for _, pid := range strings.Fields(string(out)) {
		fmt.Printf("Stopping PID %s...\n", pid)
		exec.Command("kill", "-9", pid).Run()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run killport.go <env-file>")
		os.Exit(1)
	}
	envFile := os.Args[1]

	port, err := readPort(envFile)
	if err != nil {
		fmt.Println("Warning:", err, "- skipping port kill")
		os.Exit(0)
	}
	if _, err := strconv.Atoi(port); err != nil {
		fmt.Println("Invalid port value:", port, "- skipping")
		os.Exit(0)
	}

	if runtime.GOOS == "windows" {
		killWindows(port)
	} else {
		killUnix(port)
	}
}