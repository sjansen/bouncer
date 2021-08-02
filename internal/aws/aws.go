package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// AWS provides configuration for aws-sdk clients.
type AWS struct {
	aws.Config
}

// New uses the aws-sdk prepare default client configs.
func New(ctx context.Context) (*AWS, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &AWS{
		Config: cfg,
	}, nil
}
