package service

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type StorageService struct {
	S3Client  *s3.Client
	SQSClient *sqs.Client
	Bucket    string
	QueueURL  string
}

type SQSMessage struct {
	VideoID string `json:"video_id"`
	Email   string `json:"email"`
}

func NewStorageService(s3Client *s3.Client, sqsClient *sqs.Client, bucket, queueURL string) *StorageService {
	return &StorageService{
		S3Client:  s3Client,
		SQSClient: sqsClient,
		Bucket:    bucket,
		QueueURL:  queueURL,
	}
}

func (s *StorageService) UploadFile(file multipart.File, filename string) error {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String("uploads/" + filename),
		Body:   file,
	})
	return err
}

func (s *StorageService) SendMessage(videoID, email string) error {
	payload := SQSMessage{
		VideoID: videoID,
		Email:   email,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = s.SQSClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (s *StorageService) GeneratePresignedURL(key string) (string, error) {
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

func (s *StorageService) GetBucketName() string {
    return s.Bucket
}