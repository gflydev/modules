package storagecs3

import (
	"fmt"
	"github.com/gflydev/storage/cs3"
	"net/url"
	"strings"

	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string
	fs := cs3.New()

	preSigner := PreSigner{
		MinioClient: fs.S3Client,
	}

	tempObjectKey := fmt.Sprintf("%s/%s", utils.Getenv("CS_TEMP", ""), objectKey)

	object, err := preSigner.PutObject(utils.Getenv("CS_BUCKET", ""), tempObjectKey, 60*3)
	if err != nil {
		return "", "", err
	}

	// Parse file URL
	u, _ := url.Parse(object.URL)
	bucket := utils.Getenv("CS_BUCKET", "")
	bucketCode := utils.Getenv("CS_BUCKET_CODE", "")

	preSignURL = object.URL
	fileURL = fmt.Sprintf("%s://%s/%s:%s/%s", u.Scheme, u.Host, bucketCode, bucket, tempObjectKey)

	return preSignURL, fileURL, nil
}

// LegitimizeFiles make file list available
func LegitimizeFiles(files []dto.LegitimizeItem) []dto.LegitimizeItem {
	var legitimizeItems []dto.LegitimizeItem
	fs := cs3.New()

	bucket := utils.Getenv("CS_BUCKET", "")
	bucketCode := utils.Getenv("CS_BUCKET_CODE", "")
	bucketPath := fmt.Sprintf("%s:%s/", bucketCode, bucket)

	for _, file := range files {
		object, _ := utils.RequestPath(file.File)
		object = object[1:]                                  // Remove first slash
		object = strings.Replace(object, bucketPath, "", -1) // Remove bucket path

		newObject := fmt.Sprintf("%s/%s", file.Dir, file.Name)

		fs.MakeDir(file.Dir)
		fs.Move(object, newObject)

		file.LegitimizeURL = fs.Url(newObject)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}
