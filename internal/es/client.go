package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"personalKnowledgeSearchEngine/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
)

const (
	indexName = "notes"
)

type ESSearchResponse struct {
	Hits struct {
		Hits []struct {
			ID     string      `json:"_id"`
			Source models.Note `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type ESClient struct {
	client *elasticsearch.Client
}

func NewESClient(url string) (*ESClient, error) {
	password := os.Getenv("ES_PASSWORD")
	if password == "" {
		return nil, errors.New("ES_PASSWORD environment variable not set")
	}

	cfg := elasticsearch.Config{
		Username: "elastic",
		Password: password,
		Addresses: []string{
			url,
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating es client: %s", err.Error()))
	}

	// Ping the cluster to verify the connection
	_, err = es.Info()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting response: %s", err.Error()))
	}

	return &ESClient{
		client: es,
	}, nil
}

func (es *ESClient) IndexNote(ctx context.Context, note *models.Note) error {
	noteJson, err := json.Marshal(note)
	if err != nil {
		return errors.New(fmt.Sprintf("Error marshalling note: %s", err.Error()))
	}

	_, err = es.client.Index(
		indexName,
		bytes.NewReader(noteJson),
		es.client.Index.WithContext(ctx))
	if err != nil {
		return errors.New(fmt.Sprintf("Error indexing note: %s", err.Error()))
	}

	return nil
}

func (es *ESClient) SearchNotes(query *models.SearchQuery) ([]models.Note, error) {
	queryJson, err := json.Marshal(query)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error marshalling query: %s", err.Error()))
	}

	resp, err := es.client.Search(
		es.client.Search.WithIndex(indexName),
		es.client.Search.WithBody(bytes.NewReader(queryJson)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting response: %s", err.Error()))
	}
	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("Error getting response: %s", resp.Status()))
	}

	fmt.Println(resp.String())

	var esResponse *ESSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&esResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error unmarshalling response: %s", err.Error()))
	}

	notes := make([]models.Note, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		note := hit.Source
		notes[i] = note
	}

	return notes, nil
}
