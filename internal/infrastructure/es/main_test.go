package es

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var (
	client       *elasticsearch.Client
	typedClient  *elasticsearch.TypedClient
	dao          database.DAO
	retryBackoff = backoff.NewExponentialBackOff()
)

func TestMain(m *testing.M) {
	var err error

	certPath := fmt.Sprintf("%s/ca.crt", utils.GetBasePath())
	cert, err := os.ReadFile(certPath)
	if err != nil {
		panic("Failed to read CA certificate")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "my12supers3cur3p4ssw0rd",
		CACert:   cert,

		// Retry on 429 TooManyRequests statuses
		//
		RetryOnStatus: []int{502, 503, 504, 429},

		// Configure the backoff function
		//
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	}
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	dao, err = database.NewPostgresDAO("postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable")
	if err != nil {
		panic(err)
	}

	esClient, err := NewElasticSearchClient(ElasticSearchClientParams{
		Addresses:  types.Ptr("https://localhost:9200"),
		Username:   types.Ptr("elastic"),
		Password:   types.Ptr("my12supers3cur3p4ssw0rd"),
		CACertPath: types.Ptr(certPath),
	})
	if err != nil {
		panic(err)
	}
	typedClient = esClient.typedClient

	elasticsearch.NewDefaultClient()
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	// Send the request using the client
	res, err := client.Ping()
	require.NoError(t, err)
	defer res.Body.Close()

	// Print the response status and body
	require.False(t, res.IsError())
	t.Log(res.String())
}

func TestInfo(t *testing.T) {
	// Create the requesst
	req := esapi.InfoRequest{}

	// Send the request using the typed client
	res, err := req.Do(context.Background(), client.Transport)
	require.NoError(t, err)
	defer res.Body.Close()

	// Print the response status and body
	require.False(t, res.IsError())
	t.Log(res.String())
}

func TestCreateMappingIndices(t *testing.T) {
	CreateMappingIndices(client)
}
