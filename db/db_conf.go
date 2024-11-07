package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// DBSecret represents the structure of your secret in AWS Secrets Manager
type DBSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

// getSecret retrieves and unmarshals the secret from AWS Secrets Manager
func getSecret(secretName string) (*DBSecret, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	// Get the secret value
	result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get secret value: %v", err)
	}

	// Parse the secret JSON
	var secret DBSecret
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		return nil, fmt.Errorf("unable to parse secret JSON: %v", err)
	}

	return &secret, nil
}
