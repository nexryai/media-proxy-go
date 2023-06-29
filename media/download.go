package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func fetchImage(url string) (image.Image, error) {
	core.MsgDebug(fmt.Sprintf("Donwload image: %s", url))
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		return nil, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("request ok.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
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
