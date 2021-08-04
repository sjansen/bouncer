package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSMClient provides a convenient interface to Amazon Simple Systems Manager (SSM) Parameter Store.
type SSMClient struct {
	Prefix string
	svc    *ssm.Client
}

// NewSSMClient creates a new aws-sdk S3 client.
func (aws *AWS) NewSSMClient(prefix string) (*SSMClient, error) {
	return &SSMClient{
		Prefix: prefix,
		svc:    ssm.NewFromConfig(aws.Config),
	}, nil
}

// DescribeParameters fetches multiple parameters last modified timestamp.
func (c *SSMClient) DescribeParameters(ctx context.Context, names ...string) (map[string]time.Time, error) {
	params := make([]string, 0, len(names))
	for _, name := range names {
		params = append(params, c.Prefix+name)
	}

	resp, err := c.svc.DescribeParameters(ctx, &ssm.DescribeParametersInput{
		ParameterFilters: []types.ParameterStringFilter{{
			Key:    aws.String("Name"),
			Values: params,
		}},
	})
	if err != nil {
		return nil, err
	}

	values := make(map[string]time.Time, len(names))
	for _, param := range resp.Parameters {
		name := strings.TrimPrefix(*param.Name, c.Prefix)
		values[name] = *param.LastModifiedDate
	}

	return values, nil
}

func (c *SSMClient) GetParameter(ctx context.Context, name string) (string, error) {
	resp, err := c.svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(c.Prefix + name),
		WithDecryption: true,
	})
	if err != nil {
		return "", err
	}

	return *resp.Parameter.Value, nil
}

// GetParameters fetches multiple parameters.
func (c *SSMClient) GetParameters(ctx context.Context, names ...string) (map[string]string, error) {
	params := make([]string, 0, len(names))
	for _, name := range names {
		params = append(params, c.Prefix+name)
	}

	resp, err := c.svc.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          params,
		WithDecryption: true,
	})
	if err != nil {
		return nil, err
	}

	values := make(map[string]string, len(names))
	for _, param := range resp.Parameters {
		name := strings.TrimPrefix(*param.Name, c.Prefix)
		values[name] = *param.Value
	}

	return values, nil
}

// PutParameter adds or replaces a parameter.
func (c *SSMClient) PutParameter(ctx context.Context, name string, value string, encrypt bool) error {
	name = c.Prefix + name
	input := &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: true,
	}
	if encrypt {
		input.Type = types.ParameterTypeSecureString
	} else {
		input.Type = types.ParameterTypeString
	}

	_, err := c.svc.PutParameter(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
