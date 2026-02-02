package service

import (
	"context"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"time"
)

type AWSService struct {
	S3Client  *s3.Client
	SQSClient *sqs.Client
	Bucket    string
	QueueURL  string
}

func NewAWSService() (*AWSService, error) {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsQueueURL := os.Getenv("AWS_QUEUE_URL")
	awsBucket := os.Getenv("AWS_BUCKET")
	
    if awsEndpoint == "" {
        awsEndpoint = "http://localhost:4566"
    }
	if awsAccessKeyID == "" {
        awsAccessKeyID = "teste"
    }
	if awsSecretAccessKey == "" {
        awsSecretAccessKey = "teste"
    }
	if awsRegion == "" {
        awsRegion = "us-east-1"
    }
	if awsQueueURL == "" {
        awsQueueURL = "http://localhost:4566/000000000000/video-processing-queue"
    }
	if awsBucket == "" {
        awsBucket = "fiap-x-videos"
    }	

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(awsRegion),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{URL: awsEndpoint}, nil
            })),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: awsAccessKeyID, SecretAccessKey: awsSecretAccessKey}, nil
		})),
	)
	if err != nil {
		return nil, err
	}

	return &AWSService{
		S3Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
		SQSClient: sqs.NewFromConfig(cfg),
		Bucket:    awsBucket, 
		QueueURL:  awsQueueURL,
	}, nil
}

func (s *AWSService) UploadFile(file multipart.File, filename string) error {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String("uploads/" + filename),
		Body:   file,
	})
	return err
}

func (s *AWSService) SendMessage(videoID string) error {
	_, err := s.SQSClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(videoID),
	})
	return err
}

func (s *AWSService) GeneratePresignedURL(key string) (string, error) {
    presignClient := s3.NewPresignClient(s.S3Client)
    
    req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(s.Bucket),
        Key:    aws.String(key),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = 15 * time.Minute
    })

    if err != nil {
        return "", err
    }

	//return req.URL, nil
	// força urlStr para rodar no docker compose
	urlStr := req.URL
	urlStr = strings.Replace(urlStr, "http://localstack:4566", "http://localhost:4566", 1)

    return urlStr, nil
	
}

func (s *AWSService) GetBucketName() string {
    return s.Bucket
}