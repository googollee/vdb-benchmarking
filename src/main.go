package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"vdbbench/src/client"

	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	client, err := client.NewWeaviateClient("localhost:8080")
	if err != nil {
		panic(err)
	}

	classObj := &models.Class{
		Class: "Question",
	}

	ctx := context.Background()

	exists, err := client.Schema().ClassExistenceChecker().WithClassName("Question").Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		// add the schema
		if err = client.Schema().ClassCreator().WithClass(classObj).Do(ctx); err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Open("./data/sources.json")
	if err != nil {
		log.Fatalf("open error: %v", err)
	}
	defer f.Close()

	buff := bufio.NewReader(f)

	decoder := json.NewDecoder(buff)
	for i := 0; ; i++ {
		var line []float32
		if err := decoder.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read error: %v", err)
		}

		object := models.Object{
			Class: "Question",
			Properties: map[string]any{
				"line": fmt.Sprintf("line %d", i),
			},
			Vector: line,
		}

		batchRes, err := client.Batch().ObjectsBatcher().WithObjects(&object).Do(ctx)
		if err != nil {
			log.Fatal(err)
		}
		for _, res := range batchRes {
			if res.Result.Errors != nil {
				panic(res.Result.Errors.Error)
			}
		}
	}

}
