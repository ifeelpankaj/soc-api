//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func isRunnableEnvFile(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	if name != ".env" && !strings.HasPrefix(name, ".env.") {
		return false
	}
	if strings.Contains(name, "example") || strings.Contains(name, "template") || strings.Contains(name, "sample") {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findEnvFiles() ([]string, error) {
	matches, err := filepath.Glob(".env*")
	if err != nil {
		return nil, err
	}

	envFiles := make([]string, 0, len(matches))
	for _, match := range matches {
		if isRunnableEnvFile(match) {
			envFiles = append(envFiles, filepath.ToSlash(match))
		}
	}
	sort.Strings(envFiles)
	return envFiles, nil
}

func main() {
	envFiles, err := findEnvFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to scan environment files: %v\n", err)
		os.Exit(1)
	}
	if len(envFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No runnable .env files found.")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Select environment file:")
	fmt.Fprintln(os.Stderr)
	for i, envFile := range envFiles {
		fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, envFile)
	}
	fmt.Fprintf(os.Stderr, "\nChoose [1-%d]: ", len(envFiles))

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "No selection provided.")
		os.Exit(1)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read selection: %v\n", err)
		os.Exit(1)
	}

	choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || choice < 1 || choice > len(envFiles) {
		fmt.Fprintf(os.Stderr, "Invalid selection. Choose a number from 1 to %d.\n", len(envFiles))
		os.Exit(1)
	}

	fmt.Println(envFiles[choice-1])
}
