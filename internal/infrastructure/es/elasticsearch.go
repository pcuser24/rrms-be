package es

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearchClient struct {
	client      *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

type INDICES string

const (
	LISTINGINDEX INDICES = "listings"
)

type ElasticSearchClientParams struct {
	Addresses              *string // A comma-separated list of the Elasticsearch nodes to use.
	Username               *string // Username for HTTP Basic Authentication.
	Password               *string // Password for HTTP Basic Authentication.
	CACertPath             *string // Path to the PEM file containing the CA certificate.
	Url                    *string // URL of the Elasticsearch node.
	CloudID                *string // Cloud ID from Elastic Cloud.
	APIKey                 *string // API Key for authorization.
	CertificateFingerprint *string // SHA-256 fingerprint of the certificate.
}

// NewElasticSearch creates a new ElasticSearch client instance and test the connection by sending an Ping request
func NewElasticSearchClient(params ElasticSearchClientParams) (*ElasticSearchClient, error) {
	var (
		res ElasticSearchClient
		cfg elasticsearch.Config
		err error
	)

	if params.Url != nil {
		res.client, err = elasticsearch.NewDefaultClient()
		if err != nil {
			return nil, err
		}
	}

	if params.Addresses != nil {
		addresses := strings.Split(*params.Addresses, ",")
		cfg.Addresses = addresses
	}
	if params.Username != nil {
		cfg.Username = *params.Username
	}
	if params.Password != nil {
		cfg.Password = *params.Password
	}
	if params.CACertPath != nil {
		cert, err := os.ReadFile(*params.CACertPath)
		if err != nil {
			return nil, err
		}
		cfg.CACert = cert
	}
	if params.CloudID != nil {
		cfg.CloudID = *params.CloudID
	}
	if params.APIKey != nil {
		cfg.APIKey = *params.APIKey
	}
	if params.CertificateFingerprint != nil {
		cfg.CertificateFingerprint = *params.CertificateFingerprint
	}

	// Create the client
	res.client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	res.typedClient, err = elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	// Perform a ping
	_, err = res.client.Ping()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (es *ElasticSearchClient) GetClient() *elasticsearch.Client {
	return es.client
}

func (es *ElasticSearchClient) GetTypedClient() *elasticsearch.TypedClient {
	return es.typedClient
}
