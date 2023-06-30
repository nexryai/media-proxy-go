package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/security"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
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
	case "image/webp":
		return true
	default:
		// pngはapngがあるのでfalse
		return false
	}
}

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) []byte {

	var img image.Image

	imageBuffer, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil
	}

	// 何回も参照できるようにコピー
	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	imageBuffer.Reset()

	core.MsgDebug("Content-Type: " + contentType)

	if contentType == "image/gif" {
		imgBytes, err := encodeAnimatedGifImage(bytes.NewReader(fetchedImage), contentType)

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to convert gif to webp: %v", err))
			return nil
		} else {
			return imgBytes
		}

	} else if isStaticFormat(contentType) || (contentType == "image/png" && !isAnimatedPNG(bytes.NewReader(fetchedImage))) || isStatic {

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

		imgBytes, err := ioutil.ReadAll(bytes.NewReader(fetchedImage))

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to proxy media: %v", err))
			return nil
		} else {
			return imgBytes
		}

	} else {
		//どれにも当てはまらないならnilを返してクライアントに400を返す
		return nil
	}
}
