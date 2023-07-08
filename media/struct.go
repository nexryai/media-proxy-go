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
	targetFormat   string
	isAnimated     bool
}
