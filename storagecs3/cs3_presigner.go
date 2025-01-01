package storages3

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// PreSigner encapsulates the Amazon Simple Storage Service (Amazon S3) preSigned actions
// used in the examples.
// It contains PreSignClient, a client that is used to preSigned requests to Amazon S3.
// PreSigned requests contain temporary credentials and can be made from any HTTP client.
type PreSigner struct {
	PreSignClient *s3.PresignClient
}

// GetObject makes a preSigned request that can be used to get an object from a bucket.
// The preSigned request is valid for the specified number of seconds.
func (preSigner PreSigner) GetObject(
	bucketName, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := preSigner.PreSignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// PutObject makes a preSigned request that can be used to put an object in a bucket.
// The preSigned request is valid for the specified number of seconds.
func (preSigner PreSigner) PutObject(
	bucketName, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := preSigner.PreSignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

// DeleteObject makes a preSigned request that can be used to delete an object from a bucket.
func (preSigner PreSigner) DeleteObject(bucketName, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := preSigner.PreSignClient.PresignDeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
	}
	return request, err
}
