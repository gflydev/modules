package storage

import (
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
	_utils "github.com/gflydev/modules/storage/utils"
	"github.com/gflydev/storage"
	"github.com/gflydev/storage/local"
	"time"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string

	tempObjectKey := fmt.Sprintf("%s/%s", core.TempDir, objectKey)

	preSignURL = preSignerObject(tempObjectKey)
	fileKey, _ := _utils.RequestParam(preSignURL, "G-Key")
	fileURL = fmt.Sprintf("%s/storage/tmp/%s.%s",
		core.AppURL,
		fileKey,
		utils.FileExt(objectKey),
	)

	return preSignURL, fileURL, nil
}

// LegitimizeFiles make file list available
func LegitimizeFiles(files []dto.LegitimizeItem) []dto.LegitimizeItem {
	var legitimizeItems []dto.LegitimizeItem

	for _, file := range files {
		object, _ := _utils.RequestPath(file.File)
		object = object[1:] // Remove first slash

		dir := fmt.Sprintf("%s/%s", core.AppDir, file.Dir)
		newObject := fmt.Sprintf("%s/%s", dir, file.Name)
		newObjectPath := fmt.Sprintf("%s/%s/%s", core.StorageDir, file.Dir, file.Name)

		fs := storage.Instance(local.Type)

		fs.MakeDir(dir) // Try to create new dir if not existed
		fs.Move(object, newObject)

		file.LegitimizeURL = fs.Url(newObjectPath)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}

// PreSignerObject generate Pre sign URL for a object for uploading
func preSignerObject(object string) string {
	uploadEndpoint := utils.Getenv("STORAGE_PRESIGNED_URL", "/api/v1/storage/uploads")
	// Make random data
	currentTime := time.Now().Format("20060102150405")
	randomNum := utils.RandInt64(20)
	// Token
	value := utils.Sha256(object, currentTime, randomNum)
	// File name
	fileName := fmt.Sprintf("%s.%s", value, utils.FileExt(object))

	// Caching Key
	key := fmt.Sprintf("storage:%s", value)

	// Save refresh token to Redis.
	if err := cache.Set(key, value, time.Duration(30)*time.Minute); err != nil {
		log.Fatalf("Signin error '%v'", err)
	}

	return fmt.Sprintf("%s/%s?G-Key=%s&G-Time=%s", uploadEndpoint, fileName, value, currentTime)
}
