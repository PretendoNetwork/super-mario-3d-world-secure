package main

import (
	"context"
	"os"

	"github.com/PretendoNetwork/plogger-go"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/database"
	"github.com/PretendoNetwork/super-mario-3d-world-secure/globals"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

var logger = plogger.NewLogger()

func init() {
	err := godotenv.Load()

	if err != nil {
		logger.Warning("Error loading .env file")
	}

	s3Endpoint := os.Getenv("PN_SM3DW_CONFIG_S3_ENDPOINT")
	s3Region := os.Getenv("PN_SM3DW_CONFIG_S3_REGION")
	s3AccessKey := os.Getenv("PN_SM3DW_CONFIG_S3_ACCESS_KEY")
	s3AccessSecret := os.Getenv("PN_SM3DW_CONFIG_S3_ACCESS_SECRET")

	staticCredentials := credentials.NewStaticCredentialsProvider(s3AccessKey, s3AccessSecret, "")

	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s3Endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(staticCredentials),
		config.WithEndpointResolverWithOptions(endpointResolver),
	)

	if err != nil {
		panic(err)
	}

	globals.S3Client = s3.NewFromConfig(cfg)
	globals.S3PresignClient = globals.NewPresignClient(cfg)

	database.ConnectAll()
}
