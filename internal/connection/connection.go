package connection

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

const (
	_envRegion    = "AWS_DEFAULT_REGION"
	_envAccessKey = "AWS_ACCESS_KEY_ID"
	_envSecret    = "AWS_SECRET_ACCESS_KEY"
)

type awsCredentials struct {
	// Region is the AWS region in which services are located.
	Region string

	// AccessKeyID is the nominal AWS user access key used to instantiate the SDK.
	AccessKeyID string

	// SecretKey is the nominal AWS user secret key used to instantiate the SDK.
	SecretKey string
}

func Connect() *dynamodb.Client {
	cred := awsCredentials{
		Region:      os.Getenv(_envRegion),
		AccessKeyID: os.Getenv(_envAccessKey),
		SecretKey:   os.Getenv(_envSecret),
	}

	client, err := instantiateDynamoDBClient(context.Background(), cred)
	if err != nil {
		panic(err)
	}

	return client
}

func instantiateDynamoDBClient(ctx context.Context, creds awsCredentials) (*dynamodb.Client, error) {
	dynConfig, err := configFromCredentials(ctx, creds)
	if err != nil {
		return nil, fmt.Errorf("error instantiating AWS SDK config: %w", err)
	}

	return dynamodb.NewFromConfig(dynConfig), nil
}

func configFromCredentials(ctx context.Context, creds awsCredentials) (aws.Config, error) {
	if creds.SecretKey == "" || creds.AccessKeyID == "" {
		return aws.Config{}, fmt.Errorf("could not load aws config, empty creds")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(creds.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(creds.AccessKeyID, creds.SecretKey, ""),
		),
	)

	if err != nil {
		return aws.Config{}, fmt.Errorf("could not load aws config: %w", err)
	}
	return cfg, nil

}
