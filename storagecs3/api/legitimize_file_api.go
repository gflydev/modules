package api

import (
	"github.com/gflydev/core"
	_ "github.com/gflydev/http"
	"github.com/gflydev/modules/storage/dto"
	"github.com/gflydev/modules/storage/request"
	"github.com/gflydev/modules/storagecs3"
	"github.com/gflydev/validation"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewLegitimizeFileApi As a constructor to legitimize file API.
// Related with LegitimizeFileApi
func NewLegitimizeFileApi() *LegitimizeFileApi {
	return &LegitimizeFileApi{}
}

// LegitimizeFileApi API struct.
type LegitimizeFileApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *LegitimizeFileApi) Validate(c *core.Ctx) error {
	// Parse body data.
	var legitimizeFile request.LegitimizeFile
	err := c.ParseBody(&legitimizeFile)
	if err != nil {
		c.Status(core.StatusBadRequest)

		return err
	}

	legitimizeFileDto := legitimizeFile.ToDto()

	// Validate login form.
	if errorData, err := validation.Check(legitimizeFileDto); err != nil {
		_ = c.Error(errorData)

		return err
	}

	// Store data in context.
	c.SetData(data, legitimizeFileDto)

	return nil
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Legitimize uploaded file
// @Description Legitimize uploaded file
// @Tags Storage
// @Accept json
// @Produce json
// @Success 200 {array} dto.LegitimizeItem
// @Param data body request.LegitimizeFile true "Legitimize uploaded file payload"
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Security ApiKeyAuth
// @Router /storage/legitimize-files [put]
func (h *LegitimizeFileApi) Handle(c *core.Ctx) error {
	legitimizeFile := c.GetData(data).(dto.LegitimizeFile)
	legitimizeFiles := storagecs3.LegitimizeFiles(legitimizeFile.Files)

	return c.JSON(legitimizeFiles)
}
