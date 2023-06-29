package media

import (
	"bufio"
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func convertSvgToWebp(resp *http.Response) []byte {
	w, h := 512, 512

	icon, _ := oksvg.ReadIconStream(resp.Body)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	var buf bytes.Buffer
	options, _ := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	errEncode := webp.Encode(&buf, rgba, options)
	if errEncode != nil {
		return nil
	}

	return buf.Bytes()

}

// 静止画像を取得する関数（アニメーション画像を指定しても静止画が返ってくる）
func fetchImage(url string) (image.Image, string, error) {
	core.MsgDebug(fmt.Sprintf("Download image: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		return nil, contentType, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("Request OK.")
	}

	var bodyReader io.Reader = resp.Body

	// バッファリングを行う
	bodyReader = bufio.NewReader(bodyReader)

	// SVGなら一旦webpにする
	if contentType == "image/svg+xml" {
		body := convertSvgToWebp(resp)
		if err != nil {
			return nil, contentType, fmt.Errorf("failed to convert SVG to WebP: %v", err)
		}

		bodyReader = bytes.NewReader(body)
	}

	var img image.Image
	var errDecode error

	// 適切なデコーダーを使用して画像をデコード
	switch contentType {
	case "image/webp":
		img, errDecode = webp.Decode(bodyReader, &decoder.Options{})
	default:
		img, _, errDecode = image.Decode(bodyReader)
	}

	if errDecode != nil {
		return nil, contentType, fmt.Errorf("failed to decode image: %v", errDecode)
	}

	return img, contentType, nil
}

// gif画像を取得する関数
func fetchGifImage(url string) (*gif.GIF, error) {
	core.MsgDebug(fmt.Sprintf("Donwload animated image: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download animated image. URL: " + url + ", Status: " + resp.Status)
		return nil, fmt.Errorf("failed to fetch animated image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("request ok.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gifImage, err := gif.DecodeAll(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to decode gif image: %v", err)
	}
	return gifImage, nil

}

func saveResponseToFile(resp *http.Response, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save response to file: %v", err)
	}

	return nil
}
