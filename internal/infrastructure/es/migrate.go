package es

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func createIndex(ctx context.Context, client *elasticsearch.Client, fpath string, index string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}

	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	fContent := make([]byte, fStat.Size())
	_, err = f.Read(fContent)
	if err != nil {
		return err
	}

	// Send the request using the client
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  bytes.NewReader(fContent),
	}
	res, err := req.Do(ctx, client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func CreateMappingIndices(client *elasticsearch.Client) error {
	// iterate over json files in the directory "./mappings", take the name of the file as the index name
	// and create the index with the mapping in the file
	files, err := os.ReadDir("./mappings")
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			fpath := filepath.Join("./mappings", file.Name())
			nameWithoutExtension := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			log.Println("Creating index", nameWithoutExtension, fpath)
			if err = createIndex(context.Background(), client, fpath, nameWithoutExtension); err != nil {
				return err
			}
		}
	}

	return nil
}
