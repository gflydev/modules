package api

import (
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	_ "github.com/gflydev/http"
	"github.com/gflydev/modules/storage"
	"github.com/gflydev/storage/local"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewUploadFileApi As a constructor to upload file.
// Related with UploadFileApi
func NewUploadFileApi() *UploadFileApi {
	return &UploadFileApi{}
}

// UploadFileApi API struct.
type UploadFileApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *UploadFileApi) Validate(c *core.Ctx) error {
	key := c.QueryStr("G-Key")
	fileName := c.PathVal("file_name")
	storageKey := fmt.Sprintf("storage:%s", key)

	// Reject path traversal / nested paths in the target file name before
	// it is joined with the temp directory.
	if !storage.IsSafeFileName(fileName) {
		log.Errorf("Invalid upload file name '%s'", fileName)
		return errors.New("Invalid file name")
	}

	// Check uploading key
	if _, err := cache.Get(storageKey); err != nil {
		log.Errorf("Invalid uploading token '%v'", err)
		return errors.New("Invalid uploading token")
	}

	c.SetData(data, fileName)

	return nil
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle UploadFileApi main logic for API.
// @Summary Put file to Server
// @Description Put file to Server (Local Storage). <b>Note: Don't work on Swagger 2.0</b>
// @Tags Storage
// @Accept octet-stream
// @Produce json
// @Success 204
// @Param file_name path string true "File name"
// @Param G-Key query string true "G-Key"
// @Param data body string true "Upload file"
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Security ApiKeyAuth
// @Router /storage/uploads/{file_name} [put]
func (h *UploadFileApi) Handle(c *core.Ctx) error {
	body := c.Root().PostBody()
	var fileName = c.GetData(data).(string)
	filePath := fmt.Sprintf("%s/%s", core.TempDir, fileName)

	// Get file system `local`
	fs := local.New()

	if !fs.PutData(filePath, body) {
		log.Errorf("Failed to write uploaded file '%s'", filePath)
		return errors.New("Failed to store file")
	}

	return c.NoContent()
}
