package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/security"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	"io"
	"io/ioutil"
)

func isStaticFormat(contentType string) bool {
	switch contentType {
	case "image/ico":
		// 現状icoをデコードできないのでfalseを返してデコードせずそのままプロキシするようにする
		return false
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	default:
		// pngはapngがあるのでfalse
		return false
	}
}

func isStaticImage(contentType string, bodyReader io.Reader) bool {
	if contentType == "image/png" && !isAnimatedPNG(bodyReader) {
		return true
	}
	if contentType == "image/webp" && !isAnimatedWebP(bodyReader) {
		return true
	}
	core.MsgDebug("Animated")
	return false
}

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) []byte {

	var img image.Image

	imageBuffer, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil
	}

	// FIXME: これがメモリリークの原因な気がするけどimageBufferは一度参照すると二度目以降参照できなくなる
	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	imageBuffer.Reset()

	core.MsgDebug("Content-Type: " + contentType)

	if contentType == "image/svg+xml" {
		// TODO: SVG対応
		return nil

	} else if isStaticFormat(contentType) || isStaticImage(contentType, bytes.NewReader(fetchedImage)) || isStatic {

		// 完全に静止画像のフォーマット or apngでない or static指定 ならdecodeStaticImageでデコードする
		img, err = decodeStaticImage(bytes.NewReader(fetchedImage), contentType)
		imageBuffer.Reset()

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to decode image: %v", err))
			return nil
		} else {
			core.MsgDebug("Decode ok.")
		}

		// widthLimitかheightLimitを超えている場合のみ処理する
		core.MsgDebug(fmt.Sprintf("widthLimit: %d  heightLimit: %d", widthLimit, heightLimit))

		if img.Bounds().Dx() > widthLimit || img.Bounds().Dy() > heightLimit {
			resizedImg := resizeImage(img, widthLimit, heightLimit)
			img = resizedImg
		}

		var buf bytes.Buffer

		// TODO: エンコードオプション変えられるようにする
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			return nil
		}

		errEncode := webp.Encode(&buf, img, options)
		if errEncode != nil {
			return nil
		}

		return buf.Bytes()

	} else if security.IsFileTypeBrowserSafe(contentType) {
		// どれにも当てはまらないかつブラウザセーフな形式ならそのままプロキシ
		// AVIFは敢えて無変換でプロキシする（サイズがwebpより小さくEdgeユーザーの存在を無視すれば変換する意義がほぼない）
		core.MsgDebug("Proxy image without transcode")
		return fetchedImage

	}

	//どれにも当てはまらないならnilを返してクライアントに400を返す
	return nil
}
