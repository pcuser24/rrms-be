package s3

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/random"
)

type TestS3Config struct {
	AWSRegion          string  `mapstructure:"AWS_REGION" validate:"required"`
	AWSAccessKeyID     string  `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey string  `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	AWSS3Endpoint      *string `mapstructure:"AWS_S3_ENDPOINT" validate:"omitempty"`
}

var (
	basePath = utils.GetBasePath()
	conf     TestS3Config
	s3Client *S3Client
)

func TestMain(m *testing.M) {
	viper.AddConfigPath(basePath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	v := validator.New()
	err = v.Struct(&conf)
	if err != nil {
		log.Fatalf("invalid or missing fields in config file: %v", err)
	}

	log.Println("config loaded successfully", conf, conf.AWSS3Endpoint)

	s3Client, err = NewS3Client(
		conf.AWSRegion,
		conf.AWSAccessKeyID,
		conf.AWSSecretAccessKey,
		conf.AWSS3Endpoint,
	)
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	}

	os.Exit(m.Run())
}

func createBucket(t *testing.T, bucketName string) {
	_, err := s3Client.CreateBucket(bucketName, conf.AWSRegion)
	require.NoError(t, err)

	exist, err := s3Client.BucketExists(bucketName)
	require.NoError(t, err)
	require.True(t, exist)

	t.Cleanup(func() {
		// empty the bucket
		objs, err := s3Client.ListObjects(bucketName)
		require.NoError(t, err)
		objKeys := make([]string, 0, len(objs))
		for _, obj := range objs {
			objKeys = append(objKeys, *obj.Key)
		}
		err = s3Client.DeleteObjects(bucketName, objKeys)
		require.NoError(t, err)

		_, err = s3Client.DeleteBucket(bucketName)
		require.NoError(t, err)
		exist, err := s3Client.BucketExists(bucketName)
		require.Error(t, err)
		require.False(t, exist)
	})
}

func TestCreateBucket(t *testing.T) {
	bucketName := random.RandomAlphanumericStr(5)

	createBucket(t, bucketName)
}

func TestGetPutObjectPresignedURL(t *testing.T) {
	bucketName := random.RandomAlphanumericStr(5)
	filePath, _ := filepath.Abs("./s3.go")

	createBucket(t, bucketName)

	// Open the file
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()
	fstat, err := file.Stat()
	require.NoError(t, err)

	// Get presigned URL for putting object to s3
	presignedRequest, err := s3Client.GetPutObjectPresignedURL(
		bucketName,
		"s3.go", "go", fstat.Size(),
		time.Minute,
	)
	require.NoError(t, err)
	require.NotEmpty(t, presignedRequest)

	t.Log(presignedRequest)

	// Put object to s3 using retrieved presigned URL: send PUT request which body containing the file

	// Create a PUT request with the presigned URL
	req, err := http.NewRequest(presignedRequest.Method, presignedRequest.URL, file)
	require.NoError(t, err)
	req.Header.Set("Content-Type", http.DetectContentType(make([]byte, fstat.Size())))

	// Send the request
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// List objects in the bucket
	objs, err := s3Client.ListObjects(bucketName)
	require.NoError(t, err)
	require.NotEmpty(t, objs)
	require.Equal(t, "s3.go", *objs[0].Key)
	require.Equal(t, fstat.Size(), *objs[0].Size)
}
