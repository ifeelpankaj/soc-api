//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	msgFile := flag.String("msg-file", "", "commit message file")
	flag.Parse()

	if strings.TrimSpace(*msgFile) == "" {
		fmt.Fprintln(os.Stderr, "Usage: make commit-msg-check MSG_FILE=<commit-message-file>")
		os.Exit(2)
	}

	cmd := exec.Command("go", "run", "./scripts/commitmsg", *msgFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
