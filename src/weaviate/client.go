package weaviate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

var className = "Test"

type Client struct {
	client *weaviate.Client
}

func New(addr string) (*Client, error) {
	cfg := weaviate.Config{
		Host:   addr,
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": "YOUR-OPENAI-API-KEY", // Replace with your inference API key
		},
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create weaviate client error: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Prepare(ctx context.Context) error {
	exists, err := c.client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return fmt.Errorf("weaviate: check schema error: %w", err)
	}

	if exists {
		return nil
	}

	classObj := &models.Class{
		Class: className,
	}
	if err = c.client.Schema().ClassCreator().WithClass(classObj).Do(ctx); err != nil {
		return fmt.Errorf("weaviate: create schema error: %w", err)
	}

	return nil
}

func (c *Client) Import(ctx context.Context, input *json.Decoder) error {
	for i := 0; ; i++ {
		var line []float32
		if err := input.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("weaviate: decode input error: %w", err)
		}

		object := models.Object{
			Class: className,
			Properties: map[string]any{
				"line": fmt.Sprintf("line %d", i),
			},
			Vector: line,
		}

		batchRes, err := c.client.Batch().ObjectsBatcher().WithObjects(&object).Do(ctx)
		if err != nil {
			return fmt.Errorf("weaviate: import data error: %w", err)
		}
		for _, res := range batchRes {
			if res.Result.Errors != nil {
				for _, err := range res.Result.Errors.Error {
					return fmt.Errorf("weaviate: import data error: %v", err.Message)
				}
			}
		}
	}

	return nil
}

func (c *Client) Query(ctx context.Context, input *json.Decoder) error {
	for i := 0; ; i++ {
		var line []float32
		if err := input.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("weaviate: decode query error: %w", err)
		}

		nearVector := c.client.GraphQL().NearVectorArgBuilder().WithVector(line)
		field := graphql.Field{Name: "line"}
		_additional := graphql.Field{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "certainty"}, // only supported if distance==cosine
				{Name: "distance"},  // always supported
				//{Name: "vector"},
				{Name: "id"},
			},
		}
		resp, err := c.client.GraphQL().Get().WithClassName(className).WithFields(field, _additional).WithNearVector(nearVector).WithLimit(5).Do(ctx)
		if err != nil {
			return fmt.Errorf("weaviate: query data error: %w", err)
		}
		if resp.Errors != nil {
			return fmt.Errorf("weaviate: query data error: %v", resp.Errors[0].Message)
		}
	}

	return nil
}
