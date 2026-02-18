package service

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type DbCredentials struct {
	Host     string `json:"DB_HOST"`
	Name     string `json:"DB_NAME"`
	Password string `json:"DB_PASSWORD"`
	Port     string `json:"DB_PORT"`
	Username string `json:"DB_USERNAME"`
}

func GetDatabaseSecrets(f *AWSClientFactory, secretName string) (*DbCredentials, error) {
	client := secretsmanager.NewFromConfig(f.Config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var creds DbCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &creds)
	return &creds, err
}