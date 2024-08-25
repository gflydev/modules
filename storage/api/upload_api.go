package api

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/modules/storage/response"
	"github.com/gflydev/modules/storage/transformer"
)

var (
	filename = "file"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewUploadApi As a constructor to upload file.
// Related with UploadApi
func NewUploadApi() *UploadApi {
	return &UploadApi{}
}

// UploadApi API struct.
type UploadApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *UploadApi) Validate(c *core.Ctx) error {
	// Read file header
	_, err := c.Root().FormFile(filename)
	if err != nil {
		log.Errorf("Check form file error %v", err)
		return errors.New("form file does not existed")
	}

	return nil
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle UploadApi main logic for API.
// @Summary Upload file to Server
// @Description Upload file to Server (Local Storage)
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Security ApiKeyAuth
// @Router /storage/uploads [post]
func (h *UploadApi) Handle(c *core.Ctx) error {
	var uploads []core.UploadedFile
	var err error
	var data []response.UploadedFile

	uploads, err = c.FormUpload(filename)
	if err != nil {
		log.Errorf("Upload file error %v", err)
		return errors.New("Form `%s` does not existed", filename)
	}

	for _, upload := range uploads {
		data = append(data, transformer.ToUploadedFileResponse(upload))
	}

	return c.JSON(core.Data{
		"data": data,
	})
}
