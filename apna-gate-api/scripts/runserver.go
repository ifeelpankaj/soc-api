//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func resolveEnvFile(explicitEnvFile string) (string, error) {
	envFile := strings.TrimSpace(explicitEnvFile)
	if envFile == "" {
		cmd := exec.Command("go", "run", "scripts/selectenv.go")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		var stdout bytes.Buffer
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("environment file selection failed: %w", err)
		}
		envFile = strings.TrimSpace(stdout.String())
	}

	if envFile == "" {
		return "", fmt.Errorf("no environment file selected")
	}
	if info, err := os.Stat(envFile); err != nil {
		return "", fmt.Errorf("ENV_FILE %q does not exist", envFile)
	} else if info.IsDir() {
		return "", fmt.Errorf("ENV_FILE %q is a directory", envFile)
	}

	return envFile, nil
}

func runCommand(name string, args []string, env []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}

func main() {
	modeFlag := flag.String("mode", "development", "runtime mode")
	envFileFlag := flag.String("env-file", "", "environment file to load")
	flag.Parse()

	goEnv := strings.TrimSpace(*modeFlag)
	if goEnv == "" {
		goEnv = "development"
	}

	envFile, err := resolveEnvFile(*envFileFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := runCommand("go", []string{"run", "scripts/killport.go", envFile}, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "kill-port failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting %s server with %s...\n", goEnv, envFile)
	serverEnv := append(os.Environ(), "GO_ENV="+goEnv, "ENV_FILE="+envFile)
	if err := runCommand("go", []string{"run", "cmd/server/main.go"}, serverEnv); err != nil {
		fmt.Fprintf(os.Stderr, "server exited: %v\n", err)
		os.Exit(1)
	}
}
