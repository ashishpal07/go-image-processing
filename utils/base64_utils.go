package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"net/http"
)

// Base64ToImage takes a base64 encoded string and returns an image.Image
// if it successfully decodes the string and the resulting image. Otherwise,
// it returns an error.
func Base64ToImage(base64String string) (image.Image, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, errors.New("error decoding base64 string")
	}

	img, _, err := image.Decode(bytes.NewReader(decodeBytes))
	if err != nil {
		return nil, errors.New("error decoding image")
	}

	return img, nil
}

func FetchImageFromURL(imageURL string) (image.Image, string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, "", errors.New("error while fetching image from URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("response status code is not 200 while fetching image")
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", errors.New("error decoding image after fetching from URL")
	}

	return img, format, nil
}
