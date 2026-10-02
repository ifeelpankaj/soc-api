//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Fprint(os.Stderr, "Enter migration name: ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "No migration name provided.")
		os.Exit(1)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read migration name: %v\n", err)
		os.Exit(1)
	}

	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		fmt.Fprintln(os.Stderr, "Migration name cannot be empty.")
		os.Exit(1)
	}
	name = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r + ('a' - 'A')
		case r >= '0' && r <= '9':
			return r
		case r == '-' || r == '_':
			return '_'
		default:
			return '_'
		}
	}, name)

	if err := os.MkdirAll("migrations", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create migrations directory: %v\n", err)
		os.Exit(1)
	}

	filename := filepath.Join("migrations", fmt.Sprintf("%s_%s.sql", time.Now().Format("20060102150405"), name))
	content := "-- +migrate Up\n\n-- +migrate Down\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created migration: %s\n", filepath.ToSlash(filename))
}
