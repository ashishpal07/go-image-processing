package services

import (
	"bytes"
	"errors"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

// var uploadURL = "https://api.gokapturehub.com/upload/single"

// ProcessOverlayImage takes a baseImage and an optional overlayImage and returns a byte
// representation of the resulting image. If the overlayImage is not provided, it simply
// returns a byte representation of the baseImage. If the overlayImage is provided, it resizes
// the baseImage to the size of the overlayImage, composites the overlayImage onto the baseImage
// at the origin, and returns a byte representation of the resulting image.
func ProcessOverlayImage(baseImage image.Image, overlayImage image.Image) ([]byte, error) {
	if overlayImage == nil {
		var buf bytes.Buffer
		if err := png.Encode(&buf, baseImage); err != nil {
			return nil, errors.New("error encoding image in overlay process")
		}
		return buf.Bytes(), nil
	}

	resizeBaseImage := imaging.Fill(baseImage, overlayImage.Bounds().Dx(), overlayImage.Bounds().Dy(), imaging.Center, imaging.Lanczos)
	compositeImage := imaging.Overlay(resizeBaseImage, overlayImage, image.Pt(0, 0), 1.0)

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, compositeImage); err != nil {
		return nil, errors.New("error encoding composite image in overlay process")
	}

	return buffer.Bytes(), nil
}
