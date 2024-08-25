package transformer

import (
	"github.com/gflydev/core"
	"github.com/gflydev/modules/storage/response"
)

// ToUploadedFileResponse convert gfly.UploadedFile model into response.UploadedFileResponse object.
func ToUploadedFileResponse(model core.UploadedFile) response.UploadedFile {
	return response.UploadedFile{
		Field: model.Field,
		Name:  model.Name,
		Path:  model.Path,
		Size:  model.Size,
	}
}
