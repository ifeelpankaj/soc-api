//go:build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func run(name string, args []string, env []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}

func main() {
	prod := flag.Bool("prod", false, "build optimized production binary")
	flag.Parse()

	output := "bin/server"
	if flag.NArg() > 0 {
		output = flag.Arg(0)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create build directory: %v\n", err)
		os.Exit(1)
	}

	args := []string{"build", "-o", output}
	env := os.Environ()
	if *prod {
		fmt.Println("Building production binary...")
		args = append(args, "-a", `-ldflags=-w -s`)
		env = append(env, "CGO_ENABLED=0")
	} else {
		fmt.Println("Building application...")
	}
	args = append(args, "cmd/server/main.go")

	if err := run("go", args, env); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}

	if *prod {
		fmt.Printf("Production binary created at: %s\n", output)
		return
	}
	fmt.Printf("Binary created at: %s\n", output)
}
