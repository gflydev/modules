package storagews3

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/storage/ws3"
	"net/url"
	"strings"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string
	fs := ws3.New()

	preSigner := PreSigner{
		MinioClient: fs.S3Client,
	}

	tempObjectKey := fmt.Sprintf("%s/%s", utils.Getenv("WS_TEMP", ""), objectKey)

	object, err := preSigner.PutObject(utils.Getenv("WS_BUCKET", ""), tempObjectKey, 60*3)
	if err != nil {
		return "", "", err
	}

	// Parse file URL
	u, _ := url.Parse(object.URL)
	bucket := utils.Getenv("WS_BUCKET", "")

	preSignURL = object.URL
	fileURL = fmt.Sprintf("%s://%s/%s/%s", u.Scheme, u.Host, bucket, tempObjectKey)

	return preSignURL, fileURL, nil
}

// LegitimizeFiles make file list available
func LegitimizeFiles(files []dto.LegitimizeItem) []dto.LegitimizeItem {
	var legitimizeItems []dto.LegitimizeItem
	fs := ws3.New()

	bucket := utils.Getenv("WS_BUCKET", "")
	bucketPath := fmt.Sprintf("%s/", bucket)

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
