package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	outputDir          = "test_tmp"
	coverageProfile    = "coverage.out"
	coverageSummaryOut = "coverage.txt"
	coverageHTML       = "coverage.html"
	defaultCoverageMin = 80.0
)

func main() {
	min := flag.Float64("min", -1, "minimum coverage percentage")
	enforce := flag.Bool("enforce", envBool("COVERAGE_ENFORCE"), "fail when coverage is below -min")
	flag.Parse()
	if *min < 0 {
		*min = envFloat("COVERAGE_MIN", defaultCoverageMin)
	}

	profilePath := filepath.Join(outputDir, coverageProfile)
	summaryPath := filepath.Join(outputDir, coverageSummaryOut)
	htmlPath := filepath.Join(outputDir, coverageHTML)

	fmt.Println("Running all API tests with coverage...")
	fmt.Printf("Output folder: %s\n", outputDir)
	fmt.Printf("Coverage profile: %s\n", profilePath)
	fmt.Printf("Coverage summary: %s\n", summaryPath)
	fmt.Printf("Coverage HTML: %s\n", htmlPath)

	exitOnErr(os.MkdirAll(outputDir, 0o755)) //nolint:gosec // Coverage artifacts are non-sensitive local test outputs.

	run("go", "test", "-coverprofile="+profilePath, "./...")

	summary, coverage, err := coverageSummary(profilePath)
	exitOnErr(err)
	exitOnErr(os.WriteFile(summaryPath, summary, 0o644)) //nolint:gosec // Coverage summary is non-sensitive local test output.
	fmt.Print(string(summary))

	run("go", "tool", "cover", "-html="+profilePath, "-o", htmlPath)

	fmt.Printf("Total coverage: %.1f%%\n", coverage)
	if coverage < *min {
		msg := fmt.Sprintf("coverage %.1f%% is below %.1f%%", coverage, *min)
		if *enforce {
			exitOnErr(errors.New(msg))
		}
		fmt.Printf("warning: %s; set COVERAGE_ENFORCE=1 to fail this command\n", msg)
	}
}

func coverageSummary(profile string) ([]byte, float64, error) {
	out, err := exec.CommandContext(context.Background(), "go", "tool", "cover", "-func", profile).Output() //nolint:gosec // Fixed go tool command over a locally generated coverage profile.
	if err != nil {
		return nil, 0, err
	}
	coverage, err := totalCoverage(out)
	if err != nil {
		return nil, 0, err
	}
	return out, coverage, nil
}

func totalCoverage(out []byte) (float64, error) {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		return 0, errors.New("empty coverage output")
	}
	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) == 0 {
		return 0, fmt.Errorf("invalid total coverage line: %q", lines[len(lines)-1])
	}
	raw := strings.TrimSuffix(fields[len(fields)-1], "%")
	return strconv.ParseFloat(raw, 64)
}

func run(name string, args ...string) {
	cmd := exec.CommandContext(context.Background(), name, args...) //nolint:gosec // Callers pass fixed go tool/test commands from this script.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = mustCwd()
	exitOnErr(cmd.Run())
}

func mustCwd() string {
	wd, err := os.Getwd()
	exitOnErr(err)
	return filepath.Clean(wd)
}

func envBool(name string) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return value == "1" || value == "true" || value == "yes"
}

func envFloat(name string, fallback float64) float64 {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return value
}

func exitOnErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
