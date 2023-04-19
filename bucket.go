package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

func createBucket(projectID, bucketName string) error {
	ctx := context.Background()

	// Create a storage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a new bucket handle
	bucket := client.Bucket(bucketName)

	// Set the bucket's default storage class and location
	bucketAttrs := &storage.BucketAttrs{
		StorageClass: "STANDARD",
		Location:     "US",
	}

	// Create the bucket with the specified attributes
	err = bucket.Create(ctx, projectID, bucketAttrs)
	if err != nil {
		return fmt.Errorf("Failed to create bucket: %v", err)
	}

	log.Printf("Bucket created successfully: %s", bucketName)
	return nil
}

func main() {
	// Replace with your own project ID and bucket name
	projectID := "develop-375210"
	bucketName := "my-bucket-demo-saba"

	err := createBucket(projectID, bucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	fmt.Println("Bucket creation completed successfully!")
}
