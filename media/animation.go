package media

import (
	"fmt"
	"image"
)

func extractFirstFrame(frames []image.Image) (image.Image, error) {
	if len(frames) > 0 {
		return frames[0], nil
	}

	return nil, fmt.Errorf("no frames found")
}
