package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sjansen/bouncer/internal/build"
	"github.com/sjansen/bouncer/internal/web/server"
)

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		fmt.Fprintln(os.Stderr, "This executable should be run on AWS Lambda.")
		os.Exit(1)
	}

	fmt.Println("GitSHA:", build.GitSHA)
	fmt.Println("Timestamp:", build.Timestamp)

	/*
		vars := os.Environ()
		sort.Strings(vars)
		for _, v := range vars {
			x := strings.SplitN(v, "=", 2)
			switch x[0] {
			case "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN":
				continue
			default:
				fmt.Printf("%s=%#v\n", x[0], x[1])
			}
		}
	*/

	fmt.Println("Starting server...")
	ctx := context.Background()
	s, err := server.New(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	s.StartLambdaHandler()
}
