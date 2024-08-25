package api

import (
	"github.com/gflydev/core"
	"github.com/gflydev/modules/storage/response"
	"github.com/gflydev/modules/storages3"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewPresignedURLApi As a constructor to get pre-signed URL.
// Related with NewPresignedURLApi
func NewPresignedURLApi() *PresignedURLApi {
	return &PresignedURLApi{}
}

// PresignedURLApi API struct.
type PresignedURLApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *PresignedURLApi) Validate(c *core.Ctx) error {
	filename := c.QueryStr("filename")

	// Store data in context.
	c.SetData(data, filename)

	return nil
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Get Pre-sign URL to upload
// @Description Get Pre-sign URL to upload
// @Tags Storage
// @Accept json
// @Produce json
// @Param filename query string true "Filename"
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 200 {object} response.PresignedURL
// @Security ApiKeyAuth
// @Router /storage/presigned-url [get]
func (h *PresignedURLApi) Handle(c *core.Ctx) error {
	filename := c.GetData(data).(string)

	preSignedURL, fileURL, err := storages3.PresignedURL(filename)
	if err != nil {
		return err
	}

	return c.JSON(&response.PresignedURL{
		UploadURL: preSignedURL,
		FileURL:   fileURL,
	})
}
