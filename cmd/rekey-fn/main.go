package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/sjansen/bouncer/internal/build"
	"github.com/sjansen/bouncer/internal/config"
	"github.com/sjansen/bouncer/internal/keyring"
)

type handler struct {
	keyring *keyring.KeyRing
}

func (h *handler) run(ctx context.Context) error {
	fmt.Println("Rotating keys...")
	return h.keyring.RotateJWKs(ctx)
}

func main() {
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") == "" {
		fmt.Fprintln(os.Stderr, "This executable should be run on AWS Lambda.")
		os.Exit(1)
	}

	fmt.Println("GitSHA:", build.GitSHA)
	fmt.Println("Timestamp:", build.Timestamp)
	ctx := context.Background()
	client, err := config.NewClient(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	h := &handler{
		keyring: keyring.New(client),
	}
	lambda.Start(h.run)
}
