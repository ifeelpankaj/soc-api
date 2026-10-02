//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	envFileFlag := flag.String("env-file", "", "environment file to load")
	flag.Parse()

	envFile := strings.TrimSpace(*envFileFlag)
	if envFile != "" {
		if info, err := os.Stat(envFile); err != nil {
			fmt.Fprintf(os.Stderr, "ENV_FILE %q does not exist\n", envFile)
			os.Exit(1)
		} else if info.IsDir() {
			fmt.Fprintf(os.Stderr, "ENV_FILE %q is a directory\n", envFile)
			os.Exit(1)
		}
	}

	config := ".air.linux.toml"
	if runtime.GOOS == "windows" {
		config = ".air.windows.toml"
	}

	fmt.Println("Starting with live reload...")
	fmt.Println("Make sure 'air' is installed: go install github.com/cosmtrek/air@latest")

	env := append(os.Environ(), "GO_ENV=development")
	if envFile != "" {
		env = append(env, "ENV_FILE="+envFile)
	}

	cmd := exec.Command("air", "-c", config)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "air exited: %v\n", err)
		os.Exit(1)
	}
}
