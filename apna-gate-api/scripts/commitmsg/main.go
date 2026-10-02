// Command commitmsg validates Conventional Commit messages for staged changes.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var conventionalAPISubject = regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)\(api\)!?: [a-z0-9].{9,}$`)

const gitCommandTimeout = 5 * time.Second

func main() {
	if len(os.Args) != 2 {
		exitWithUsage()
	}

	hasAPIChanges, err := stagedAPIChanges()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to inspect staged API files: %v\n", err)
		os.Exit(1)
	}
	if !hasAPIChanges {
		fmt.Println("No staged files; skipping API commit message check.")
		return
	}

	message, err := readCommitMessage(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read commit message: %v\n", err)
		os.Exit(1)
	}

	subject := firstSubjectLine(string(message))
	if conventionalAPISubject.MatchString(subject) && !isVagueSubject(subject) {
		return
	}

	fmt.Fprintln(os.Stderr, "API commits must use Conventional Commits with the api scope.")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  feat(api): add visitor invite expiry")
	fmt.Fprintln(os.Stderr, "  fix(api): validate guard entry token")
	fmt.Fprintln(os.Stderr, "  test(api): cover visitor invite expiry")
	os.Exit(1)
}

func stagedAPIChanges() (bool, error) {
	root, err := gitRoot()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "-C", root, "diff", "--cached", "--name-only") // #nosec G204 -- fixed git command with discovered repository root.
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return false, fmt.Errorf("%w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return false, err
	}

	return strings.TrimSpace(string(out)) != "", nil
}

func gitRoot() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func readCommitMessage(path string) ([]byte, error) {
	message, err := os.ReadFile(path) // #nosec G304,G703 -- Git passes the commit message file path to this hook.
	if err == nil || filepath.IsAbs(path) {
		return message, err
	}

	return os.ReadFile(filepath.Join("..", path)) // #nosec G304,G703 -- fallback for Lefthook running from a nested directory.
}

func firstSubjectLine(message string) string {
	for _, line := range strings.Split(message, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line
	}
	return ""
}

func isVagueSubject(subject string) bool {
	_, description, ok := strings.Cut(subject, ": ")
	if !ok {
		return true
	}

	normalized := strings.ToLower(strings.TrimSpace(description))
	vagueDescriptions := map[string]struct{}{
		"change":           {},
		"changes":          {},
		"fix":              {},
		"fixes":            {},
		"update":           {},
		"updates":          {},
		"wip":              {},
		"work in progress": {},
	}

	_, vague := vagueDescriptions[normalized]
	return vague
}

func exitWithUsage() {
	fmt.Fprintln(os.Stderr, "usage: go run ./scripts/commitmsg <commit-message-file>")
	os.Exit(2)
}
