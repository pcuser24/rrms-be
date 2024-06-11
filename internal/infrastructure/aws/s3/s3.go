package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3Client interface {
	CreateBucket(name string, region string) (*s3.CreateBucketOutput, error)
	BucketExists(bucketName string) (bool, error)
	DeleteBucket(bucketName string) (*s3.DeleteBucketOutput, error)
	GetPutObjectPresignedURL(bucketName string, objectKey, contentType string, contentLength int64, lifetime time.Duration) (*v4.PresignedHTTPRequest, error)
	ListObjects(bucketName string) ([]types.Object, error)
	UploadLargeObject(bucketName string, objectKey string, largeObject []byte) error
	DownloadFile(bucketName string, objectKey string, fileName string) error
	DeleteObjects(bucketName string, objectKeys []string) error
}

type s3Client struct {
	s3Client  *s3.Client
	presigner *s3.PresignClient
}

// Set environment variables: AWS_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
func NewS3Client(conf *aws.Config) S3Client {
	c := s3.NewFromConfig(*conf, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &s3Client{
		s3Client:  c,
		presigner: s3.NewPresignClient(c),
	}
}

func (c *s3Client) CreateBucket(name string, region string) (*s3.CreateBucketOutput, error) {
	return c.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
}

func (c *s3Client) BucketExists(bucketName string) (bool, error) {
	_, err := c.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
			// exists = false
			// err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
		return false, err
	} else {
		// log.Printf("Bucket %v exists and you already own it.", bucketName)
		return true, nil
	}
}

// DeleteBucket deletes a bucket. The bucket must be empty or an error is returned.
func (c *s3Client) DeleteBucket(bucketName string) (*s3.DeleteBucketOutput, error) {
	return c.s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
}

// GetPutObjectPresignedURL makes a presigned URL that can be used to PUT an object in a bucket.
// The presigned URL is valid for the specified duration.
func (c *s3Client) GetPutObjectPresignedURL(
	bucketName string,
	objectKey, contentType string, contentLength int64,
	lifetime time.Duration,
) (*v4.PresignedHTTPRequest, error) {
	return c.presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = lifetime
	})
}

// ListObjects lists the objects in a bucket.
func (c *s3Client) ListObjects(bucketName string) ([]types.Object, error) {
	result, err := c.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}

// UploadLargeObject uses an upload manager to upload data to an object in a bucket.
// The upload manager breaks large data into parts and uploads the parts concurrently.
func (c *s3Client) UploadLargeObject(bucketName string, objectKey string, largeObject []byte) error {
	largeBuffer := bytes.NewReader(largeObject)
	var partMiBs int64 = 10
	uploader := manager.NewUploader(c.s3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   largeBuffer,
	})
	if err != nil {
		log.Printf("Couldn't upload large object to %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}

	return err
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (c *s3Client) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := c.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

// DeleteObjects deletes objects from a bucket.
func (c *s3Client) DeleteObjects(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	output, err := c.s3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		log.Printf("Deleted %v objects.\n", len(output.Deleted))
	}
	return err
}
