package storagews3

import (
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/storage/ws3"
	"net/url"
	"path/filepath"
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
	u, err := url.Parse(object.URL)
	if err != nil {
		return "", "", err
	}
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
		// Validate user-supplied path components to prevent traversal.
		if !isSafeRelPath(file.Dir) || !isSafeFileName(file.Name) {
			log.Errorf("Legitimize file rejected: unsafe dir '%s' or name '%s'", file.Dir, file.Name)
			continue
		}

		object, err := utils.RequestPath(file.File)
		if err != nil || len(object) == 0 {
			log.Errorf("Legitimize file rejected: bad source path '%s' (%v)", file.File, err)
			continue
		}
		object = strings.TrimPrefix(object, "/")        // Remove first slash
		object = strings.TrimPrefix(object, bucketPath) // Remove bucket path prefix

		newObject := fmt.Sprintf("%s/%s", file.Dir, file.Name)

		if !fs.MakeDir(file.Dir) {
			log.Errorf("Legitimize file: make dir '%s' failed", file.Dir)
			continue
		}
		if !fs.Move(object, newObject) {
			log.Errorf("Legitimize file: move '%s' -> '%s' failed", object, newObject)
			continue
		}

		file.LegitimizeURL = fs.Url(newObject)

		legitimizeItems = append(legitimizeItems, file)
	}

	return legitimizeItems
}

// isSafeFileName reports whether name is a plain file name with no path
// separators or traversal sequences.
func isSafeFileName(name string) bool {
	if name == "" {
		return false
	}
	if strings.ContainsAny(name, "/\\") {
		return false
	}

	return isSafeRelPath(name)
}

// isSafeRelPath reports whether p is a relative path that stays within its
// base directory (no `..` traversal, no absolute path).
func isSafeRelPath(p string) bool {
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
