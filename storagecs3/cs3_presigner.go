package storagecs3

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
)

// PreSigner encapsulates the MinIO/S3 preSigned actions
// used in the examples.
// It contains MinioClient, a client that is used to preSigned requests to S3-compatible storage.
// PreSigned requests contain temporary credentials and can be made from any HTTP client.
type PreSigner struct {
	MinioClient *minio.Client
}

// PresignedHTTPRequest represents a presigned HTTP request
type PresignedHTTPRequest struct {
	URL string
}

// GetObject makes a preSigned request that can be used to get an object from a bucket.
// The preSigned request is valid for the specified number of seconds.
func (preSigner PreSigner) GetObject(
	bucketName, objectKey string, lifetimeSecs int64) (*PresignedHTTPRequest, error) {
	expires := time.Duration(lifetimeSecs) * time.Second
	presignedURL, err := preSigner.MinioClient.PresignedGetObject(context.TODO(), bucketName, objectKey, expires, nil)
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return nil, err
	}
	return &PresignedHTTPRequest{URL: presignedURL.String()}, nil
}

// PutObject makes a preSigned request that can be used to put an object in a bucket.
// The preSigned request is valid for the specified number of seconds.
func (preSigner PreSigner) PutObject(
	bucketName, objectKey string, lifetimeSecs int64) (*PresignedHTTPRequest, error) {
	expires := time.Duration(lifetimeSecs) * time.Second
	presignedURL, err := preSigner.MinioClient.PresignedPutObject(context.TODO(), bucketName, objectKey, expires)
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return nil, err
	}
	return &PresignedHTTPRequest{URL: presignedURL.String()}, nil
}

// DeleteObject makes a preSigned request that can be used to delete an object from a bucket.
// Note: MinIO doesn't support presigned DELETE operations, so this method performs direct deletion
func (preSigner PreSigner) DeleteObject(bucketName, objectKey string) (*PresignedHTTPRequest, error) {
	// MinIO doesn't support presigned DELETE operations
	// We'll perform the delete operation directly instead
	err := preSigner.MinioClient.RemoveObject(context.TODO(), bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Couldn't delete object %v. Here's why: %v\n", objectKey, err)
		return nil, err
	}

	// Return a dummy URL since the operation was completed directly
	return &PresignedHTTPRequest{URL: "deleted"}, nil
}
