package s3

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	conf      aws.Config
	s3Client  *s3.Client
	presigner *s3.PresignClient
}

// Set environment variables: AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
func NewS3Client() (*S3Client, error) {
	conf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(conf)

	return &S3Client{
		conf:      conf,
		s3Client:  s3Client,
		presigner: s3.NewPresignClient(s3Client),
	}, nil
}

func (c *S3Client) CreateBucket(name string, region string) error {
	_, err := c.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		log.Printf("Couldn't create bucket %v in Region %v. Here's why: %v\n",
			name, region, err)
	}
	return err
}

func (c *S3Client) BucketExists(bucketName string) (bool, error) {
	_, err := c.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true

	return exists, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified duration.
func (c *S3Client) PutObject(
	bucketName string,
	objectKey, contentType string, contentLength int64,
	lifetime time.Duration,
) (*v4.PresignedHTTPRequest, error) {
	return c.presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: contentLength,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = lifetime
	})
}
