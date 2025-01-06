package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
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
