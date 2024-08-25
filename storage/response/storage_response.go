package response

// PresignedURL struct to describe pre-signed URL response.
type PresignedURL struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
}

type UploadedFile struct {
	Field string `json:"field"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
}
