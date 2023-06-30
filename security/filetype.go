package security

const (
	FileTypePNG       = "image/png"
	FileTypeGIF       = "image/gif"
	FileTypeJPEG      = "image/jpeg"
	FileTypeWebP      = "image/webp"
	FileTypeAVIF      = "image/avif"
	FileTypeAPNG      = "image/apng"
	FileTypeBMP       = "image/bmp"
	FileTypeTIFF      = "image/tiff"
	FileTypeIcon      = "image/x-icon"
	FileTypeOpus      = "audio/opus"
	FileTypeOGG       = "video/ogg"
	FileTypeOGGAudio  = "audio/ogg"
	FileTypeOGGApp    = "application/ogg"
	FileTypeQuicktime = "video/quicktime"
	FileTypeMP4       = "video/mp4"
	FileTypeMP4Audio  = "audio/mp4"
	FileTypeM4V       = "video/x-m4v"
	FileTypeM4A       = "audio/x-m4a"
	FileType3GPP      = "video/3gpp"
	FileType3GPP2     = "video/3gpp2"
	FileTypeMPEG      = "video/mpeg"
	FileTypeMPEGAudio = "audio/mpeg"
	FileTypeWebM      = "video/webm"
	FileTypeWebMAudio = "audio/webm"
	FileTypeAAC       = "audio/aac"
	FileTypeFLAC      = "audio/flac"
	FileTypeWAV       = "audio/wav"
	FileTypeOldFLAC   = "audio/x-flac"
	FileTypeWave      = "audio/vnd.wave"
)

var FileTypeBrowserSafe = []string{
	FileTypePNG,
	FileTypeGIF,
	FileTypeJPEG,
	FileTypeWebP,
	FileTypeAVIF,
	FileTypeAPNG,
	FileTypeBMP,
	FileTypeTIFF,
	FileTypeIcon,
	FileTypeOpus,
	FileTypeOGG,
	FileTypeOGGAudio,
	FileTypeOGGApp,
	FileTypeQuicktime,
	FileTypeMP4,
	FileTypeMP4Audio,
	FileTypeM4V,
	FileTypeM4A,
	FileType3GPP,
	FileType3GPP2,
	FileTypeMPEG,
	FileTypeMPEGAudio,
	FileTypeWebM,
	FileTypeWebMAudio,
	FileTypeAAC,
	FileTypeFLAC,
	FileTypeWAV,
	FileTypeOldFLAC,
	FileTypeWave,
}

func IsFileTypeBrowserSafe(fileType string) bool {
	for _, ft := range FileTypeBrowserSafe {
		if fileType == ft {
			return true
		}
	}
	return false
}
