package storages3

import (
	"fmt"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/storage/s3"
	"net/url"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string
	fs := s3.New()

	preSigner := PreSigner{
		PreSignClient: awsS3.NewPresignClient(fs.S3Client),
	}

	tempObjectKey := fmt.Sprintf("%s/%s", utils.Getenv("AWS_S3_TEMP", ""), objectKey)

	object, err := preSigner.PutObject(utils.Getenv("AWS_S3_BUCKET", ""), tempObjectKey, 60*3)
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
	fs := s3.New()

	for _, file := range files {
		object, _ := utils.RequestPath(file.File)
		object = object[1:] // Remove first slash

		newObject := fmt.Sprintf("%s/%s", file.Dir, file.Name)

		fs.Move(object, newObject)
		file.LegitimizeURL = fs.Url(newObject)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}
