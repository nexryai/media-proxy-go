package media

type ProxyOpts struct {
	Url          string
	WidthLimit   int
	HeightLimit  int
	IsStatic     bool
	TargetFormat string
}

type transcodeImageOpts struct {
	imageBufferPtr *[]byte
	widthLimit     int
	heightLimit    int
	originalFormat string
	targetFormat   string
	isAnimated     bool
}

type ffmpegOpts struct {
	imageBufferPtr *[]byte
	shouldResize   bool
	width          int
	height         int
	encoder        string
	targetFormat   string
}
