package request

import "github.com/gflydev/modules/storage/dto"

// PreSignURL struct to describe Pre-sign URL.
type PreSignURL struct {
	dto.PreSignURL
}

// ToDto Convert to PreSignURL DTO object.
func (r PreSignURL) ToDto() dto.PreSignURL {
	return dto.PreSignURL{
		Filename: r.Filename,
	}
}

// LegitimizeFile struct to describe to legitimize files.
type LegitimizeFile struct {
	dto.LegitimizeFile
}

// ToDto Convert to LegitimizeFile DTO object.
func (r LegitimizeFile) ToDto() dto.LegitimizeFile {
	return dto.LegitimizeFile{
		Files: r.Files,
	}
}
