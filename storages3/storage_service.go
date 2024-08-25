package storages3

import (
	"fmt"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gflydev/modules/storage/dto"
	_utils "github.com/gflydev/modules/storage/utils"
	"github.com/gflydev/storage"
	"github.com/gflydev/storage/s3"
	"net/url"
	"os"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string
	fsS3 := s3.New()

	preSigner := PreSigner{
		PreSignClient: awsS3.NewPresignClient(fsS3.S3Client),
	}

	tempObjectKey := fmt.Sprintf("%s/%s", os.Getenv("AWS_S3_TEMP"), objectKey)

	object, err := preSigner.PutObject(os.Getenv("AWS_S3_BUCKET"), tempObjectKey, 60*3)
	if err != nil {
		return "", "", err
	}

	// Parse file URL
	u, _ := url.Parse(object.URL)

	preSignURL = object.URL
	fileURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

	return preSignURL, fileURL, nil
}

// LegitimizeFiles make file list available
func LegitimizeFiles(files []dto.LegitimizeItem) []dto.LegitimizeItem {
	var legitimizeItems []dto.LegitimizeItem
	fsS3 := storage.Instance(s3.Type)

	for _, file := range files {
		object, _ := _utils.RequestPath(file.File)
		object = object[1:] // Remove first slash

		newObject := fmt.Sprintf("%s/%s", file.Dir, file.Name)

		fsS3.Move(object, newObject)
		file.LegitimizeURL = fsS3.Url(newObject)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}
