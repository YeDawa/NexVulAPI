package services

import (
	"io"
	"os"
	"fmt"
	"context"

	"encoding/hex"
	"crypto/sha256"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type R2 struct{}

func (r2 R2) UploadFile(file io.Reader, key string, contentType string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
				}, nil
			},
		))),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("R2_ENDPOINT"),
					SigningRegion: os.Getenv("R2_BUCKET_REGION"),
				}, nil
			},
		)),
	)

	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	if bucketName == "" {
		return fmt.Errorf("R2_BUCKET_NAME is not set")
	}

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &key,
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPrivate,
	})

	if err != nil {
		return fmt.Errorf("unable to upload file to bucket %s: %v", bucketName, err)
	}

	return nil
}

func (r2 R2) ReadFile(key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
				}, nil
			},
		))),

		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("R2_ENDPOINT"),
					SigningRegion: os.Getenv("R2_BUCKET_REGION"),
				}, nil
			},
		)),
	)

	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	resp, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})

	if err != nil {
		return "", fmt.Errorf("unable to get object %s from bucket %s, %v", key, bucketName, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read object body, %v", err)
	}

	return string(body), nil
}

func (r2 R2) EditFile(file io.Reader, key string, contentType string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
				}, nil
			},
		))),

		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("R2_ENDPOINT"),
					SigningRegion: os.Getenv("R2_BUCKET_REGION"),
				}, nil
			},
		)),
	)

	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &key,
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPrivate,
	})

	if err != nil {
		return fmt.Errorf("unable to edit object %s in bucket %s, %v", key, bucketName, err)
	}

	return nil
}

func (r2 R2) DeleteFile(key string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
				}, nil
			},
		))),

		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("R2_ENDPOINT"),
					SigningRegion: os.Getenv("R2_BUCKET_REGION"),
				}, nil
			},
		)),
	)

	if err != nil {
		return fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})

	if err != nil {
		return fmt.Errorf("unable to delete object %s from bucket %s, %v", key, bucketName, err)
	}

	return nil
}

func (r2 R2) ChecksumFile(key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
					SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
				}, nil
			},
		))),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           os.Getenv("R2_ENDPOINT"),
					SigningRegion: os.Getenv("R2_BUCKET_REGION"),
				}, nil
			},
		)),
	)

	if err != nil {
		return "", fmt.Errorf("unable to load SDK config: %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("R2_BUCKET_NAME")

	resp, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})

	if err != nil {
		return "", fmt.Errorf("unable to get object %s from bucket %s: %v", key, bucketName, err)
	}

	defer resp.Body.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, resp.Body); err != nil {
		return "", fmt.Errorf("unable to read object body: %v", err)
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	return checksum, nil
}