package media

type ProxyRequest struct {
	Url          string
	WidthLimit   int
	HeightLimit  int
	IsStatic     bool
	IsEmoji      bool
	UseAVIF      bool
	TargetFormat string
}
