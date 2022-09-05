package gonetable_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func MustLoadLocalDDBConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("local"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
				Source:          "dummy",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	return cfg
}

func DeleteTableIfExists(ctx context.Context, ddb *dynamodb.Client, tablename string) error {
	ts, err := ddb.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}
	for _, n := range ts.TableNames {
		if n == tablename {
			_, err := ddb.DeleteTable(ctx, &dynamodb.DeleteTableInput{
				TableName: aws.String(tablename),
			})
			return err
		}
	}
	return nil
}
