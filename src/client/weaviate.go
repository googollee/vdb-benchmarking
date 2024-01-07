package client

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func NewWeaviateClient(endpoint string) (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Host:   endpoint,
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": "YOUR-OPENAI-API-KEY", // Replace with your inference API key
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create weaviate client error: %w", err)
	}

	return client, nil
}
