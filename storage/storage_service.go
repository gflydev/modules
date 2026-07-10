package storage

import (
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/storage/local"
	"path/filepath"
	"strings"
	"time"
)

// PresignedURL generate pre-signed upload URL from Local storage
func PresignedURL(objectKey string) (string, string, error) {
	var preSignURL, fileURL string

	// Reject traversal / absolute paths in the caller-supplied object key.
	if !IsSafeRelPath(objectKey) {
		return "", "", fmt.Errorf("invalid object key %q", objectKey)
	}

	tempObjectKey := fmt.Sprintf("%s/%s", core.TempDir, objectKey)

	preSignURL, err := preSignerObject(tempObjectKey)
	if err != nil {
		return "", "", err
	}

	fileKey, _ := utils.RequestParam(preSignURL, "G-Key")
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
	fs := local.New()

	for _, file := range files {
		// Validate user-supplied path components to prevent traversal
		// outside of the application / storage roots.
		if !IsSafeRelPath(file.Dir) || !IsSafeFileName(file.Name) {
			log.Errorf("Legitimize file rejected: unsafe dir '%s' or name '%s'", file.Dir, file.Name)
			continue
		}

		object, err := utils.RequestPath(file.File)
		if err != nil || len(object) == 0 {
			log.Errorf("Legitimize file rejected: bad source path '%s' (%v)", file.File, err)
			continue
		}
		object = strings.TrimPrefix(object, "/") // Remove first slash

		dir := fmt.Sprintf("%s/%s", core.AppDir, file.Dir)
		newObject := fmt.Sprintf("%s/%s", dir, file.Name)
		newObjectPath := fmt.Sprintf("%s/%s/%s", core.StorageDir, file.Dir, file.Name)

		if !fs.MakeDir(dir) { // Try to create new dir if not existed
			log.Errorf("Legitimize file: make dir '%s' failed", dir)
			continue
		}
		if !fs.Move(object, newObject) {
			log.Errorf("Legitimize file: move '%s' -> '%s' failed", object, newObject)
			continue
		}

		file.LegitimizeURL = fs.Url(newObjectPath)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}

// PreSignerObject generate Pre sign URL for a object for uploading
func preSignerObject(object string) (string, error) {
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

	// Save upload token to cache. On failure return an error instead of
	// terminating the process (previously log.Fatalf → os.Exit, a remote DoS).
	if err := cache.Set(key, value, time.Duration(30)*time.Minute); err != nil {
		log.Errorf("Presigned URL cache error '%v'", err)

		return "", err
	}

	return fmt.Sprintf("%s/%s?G-Key=%s&G-Time=%s", uploadEndpoint, fileName, value, currentTime), nil
}

// IsSafeFileName reports whether name is a plain file name with no path
// separators or traversal sequences.
func IsSafeFileName(name string) bool {
	if name == "" {
		return false
	}
	if strings.ContainsAny(name, "/\\") {
		return false
	}

	return IsSafeRelPath(name)
}

// IsSafeRelPath reports whether p is a relative path that stays within its
// base directory (no `..` traversal, no absolute path).
func IsSafeRelPath(p string) bool {
	if p == "" || strings.Contains(p, "\x00") {
		return false
	}
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") {
		return false
	}

	cleaned := filepath.ToSlash(filepath.Clean(p))
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return false
	}

	return true
}
