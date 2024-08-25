package dto

// PreSignURL struct to describe PreSign URL.
type PreSignURL struct {
	Filename string `json:"filename" example:"avatar.png" validate:"required,lte=255"`
}

// LegitimizeItem struct to describe legitimize item.
type LegitimizeItem struct {
	File          string `json:"file" example:"https://902-local.s3.us-west-1.amazonaws.com/tmp/avatar.png" validate:"required,lte=255"`
	Name          string `json:"name" example:"avatar.png" validate:"required,lte=255"`
	Dir           string `json:"dir" example:"avatars" validate:"required,lte=255"`
	LegitimizeURL string `json:"legitimize_url" example:""` // A legitimize URL (public URL). It will be created from backend.
}

// LegitimizeFile struct to describe legitimize list of uploaded File.
type LegitimizeFile struct {
	Files []LegitimizeItem `json:"files" validate:"required"`
}
