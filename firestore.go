package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	
)

func main() {
	// Set up a new Firestore client with default credentials
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "develop-375210")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Create a new collection
	_, err = client.Collection("my-collection-saba").Doc("my-document").Set(ctx, map[string]interface{}{
		"field1": "value1",
		"field2": 42,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Collection created successfully!")
}
