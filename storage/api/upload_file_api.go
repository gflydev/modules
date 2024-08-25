package api

import (
	"fmt"
	"github.com/gflydev/cache"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/storage"
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
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Security ApiKeyAuth
// @Router /storage/uploads/{file_name} [put]
func (h *UploadFileApi) Handle(c *core.Ctx) error {
	body := c.Root().PostBody()
	var fileName = c.GetData(data).(string)
	filePath := fmt.Sprintf("%s/%s", core.TempDir, fileName)

	// Get file system `local`
	fs := storage.Instance(local.Type)

	fs.PutData(filePath, body)

	return c.NoContent()
}
