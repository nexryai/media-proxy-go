package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
)

func fetchImage(url string) (image.Image, error) {
	core.MsgDebug(fmt.Sprintf("Donwload image: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		core.MsgWarn("Failed to download image. url: " + url)
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		return nil, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("request ok.")
	}

	if core.IsDebugMode() {
		err = saveResponseToFile(resp, "debug.bin")
		if err != nil {
			core.MsgWarn("Failed to save debug image: " + err.Error())
			// ファイル保存のエラーは無視して続行
		}
	}

	var img image.Image

	core.MsgDebug(fmt.Sprintf("Decode image: %s", url))

	// FIXME: なんかunknown formatになる
	img, _, err = image.Decode(resp.Body)
	if err != nil {
		core.MsgWarn("Failed to decode image. url: " + url)
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
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
