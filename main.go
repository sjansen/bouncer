package main

import (
	"fmt"
	"os"

	"github.com/sjansen/bouncer/internal/cli"
)

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		fmt.Fprintln(os.Stderr, "This executable should not be run on AWS Lambda.")
		os.Exit(1)
	}

	cli.ParseAndRun()
}
