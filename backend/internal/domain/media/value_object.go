package media

type MediaMetadata struct {
	Width      int
	Height     int
	Duration   float64
	Format     string
	FileSize   int64
	Resolution string
}

type GenerationParams struct {
	Prompt         string
	NegativePrompt string
	Style          string
	Quality        string
	ReferenceImage string
	Seed           int64
	Steps          int
	CFGScale       float64
}

func NewImageMetadata(width, height int, format string, fileSize int64) MediaMetadata {
	return MediaMetadata{
		Width:      width,
		Height:     height,
		Format:     format,
		FileSize:   fileSize,
		Resolution: formatResolution(width, height),
	}
}

func NewVideoMetadata(width, height int, duration float64, format string, fileSize int64) MediaMetadata {
	return MediaMetadata{
		Width:      width,
		Height:     height,
		Duration:   duration,
		Format:     format,
		FileSize:   fileSize,
		Resolution: formatResolution(width, height),
	}
}

func formatResolution(width, height int) string {
	return string(rune(width)) + "x" + string(rune(height))
}
