package service

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type AWSClientFactory struct {
	Config aws.Config
}

func NewAWSClientFactory(ctx context.Context, region, endpoint, key, secret string) *AWSClientFactory {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if endpoint != "" {
		opts = append(opts, config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: region}, nil
			},
		)))
	}

	if key != "" && secret != "" {
		opts = append(opts, config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: key, SecretAccessKey: secret}, nil
		})))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	return &AWSClientFactory{Config: cfg}
}

func (f *AWSClientFactory) NewS3Client() *s3.Client {
	return s3.NewFromConfig(f.Config, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func (f *AWSClientFactory) NewSQSClient() *sqs.Client {
	return sqs.NewFromConfig(f.Config)
}