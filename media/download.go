package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
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

func fetchImage(url string) (image.Image, string, error) {
	core.MsgDebug(fmt.Sprintf("Donwload image: %s", url))

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
		core.MsgDebug("request ok.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, contentType, err
	}

	var img image.Image
	core.MsgDebug(contentType)

	// SVGなら一旦webpにする
	// ToDo: もっとマシな書き方あるだろうけどsvgってそんな使うもん？
	if contentType == "image/svg+xml" {
		body = convertSvgToWebp(resp)
	} else if contentType == "image/webp" {
		imgDecoded, err := webp.Decode(bytes.NewReader(body), &decoder.Options{})
		if err != nil {
			return nil, contentType, fmt.Errorf("failed to decode webp: %v", err)
		}
		img = imgDecoded
	} else {
		imgDecoded, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			return nil, contentType, fmt.Errorf("failed to decode image: %v", err)
		}
		img = imgDecoded
	}
	return img, contentType, nil
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
