package main

import (
	"context"
	"fmt"
	"vdbbench/client"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func main() {
	client, err := client.NewWeaviateClient("localhost:8080")
	if err != nil {
		panic(err)
	}

	// classObj := &models.Class{
	// Class: "Question",
	// }

	// // add the schema
	// err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	// if err != nil {
	// panic(err)
	// }

	// items := []struct {
	// Category string
	// Question string
	// Answer   string
	// }{
	// {"SCIENCE", "This organ removes excess glucose from the blood & stores it as glycogen", "Liver"},
	// {"ANIMALS", "It's the only living mammal in the order Proboseidea", "Elephant"},
	// {"SCIENCE", "In 70-degree air, a plane traveling at about 1,130 feet per second breaks it", "Sound barrier"},
	// }
	// objects := make([]*models.Object, len(items))
	// for i := range items {
	// objects[i] = &models.Object{
	// Class: "Question",
	// Properties: map[string]any{
	// "category": items[i].Category,
	// "question": items[i].Question,
	// "answer":   items[i].Answer,
	// },
	// Vector: []float32{0.1, 0.4, 0.2},
	// }
	// }

	// batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	// if err != nil {
	// panic(err)
	// }
	// for _, res := range batchRes {
	// if res.Result.Errors != nil {
	// panic(res.Result.Errors.Error)
	// }
	// }

	fields := []graphql.Field{
		{Name: "question"},
		{Name: "answer"},
		{Name: "category"},
	}

	near := client.GraphQL().
		NearVectorArgBuilder().
		WithDistance(0.2).
		WithVector([]float32{0.3})

	result, err := client.GraphQL().Get().
		WithClassName("Question").
		WithFields(fields...).
		WithNearVector(near).
		WithLimit(2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", result.Errors[0].Message)
	fmt.Printf("%+v\n", result)
}
