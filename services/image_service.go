package services

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"encoding/json"

	"github.com/disintegration/imaging"
)

// var uploadURL = "https://api.gokapturehub.com/upload/single"

// ProcessOverlayImage takes a baseImage and an optional overlayImage and returns a byte
// representation of the resulting image. If the overlayImage is not provided, it simply
// returns a byte representation of the baseImage. If the overlayImage is provided, it resizes
// the baseImage to the size of the overlayImage, composites the overlayImage onto the baseImage
// at the origin, and returns a byte representation of the resulting image.
func ProcessOverlayImage(baseImage image.Image, overlayImage image.Image) (string, error) {
	var processedImage image.Image

	if overlayImage == nil {
		processedImage = baseImage
	} else {
		resizedBaseImage := imaging.Fill(baseImage, overlayImage.Bounds().Dx(), overlayImage.Bounds().Dy(), imaging.Center, imaging.Lanczos)
		processedImage = imaging.Overlay(resizedBaseImage, overlayImage, image.Pt(0, 0), 1.0)
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, processedImage); err != nil {
		return "", errors.New("failed to encode the processed image to PNG: " + err.Error())
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "image.png")
	if err != nil {
		return "", errors.New("failed to create form file for upload: " + err.Error())
	}

	if _, err = io.Copy(part, &buffer); err != nil {
		return "", errors.New("failed to write image to form file: " + err.Error())
	}

	if err := writer.Close(); err != nil {
		return "", errors.New("failed to close multipart writer: " + err.Error())
	}

	const uploadURL = "https://api.gokapturehub.com/upload/single"

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", errors.New("failed to create the upload request: " + err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("failed to send the upload request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("upload request failed with status: " + resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read the response body: " + err.Error())
	}

	type UploadResponse struct {
		Status    string `json:"status"`
		Data      string `json:"data"`
		Timestamp string `json:"timestamp"`
		Message   string `json:"message"`
	}

	var uploadResp UploadResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return "", errors.New("failed to parse the upload response: " + err.Error())
	}

	if uploadResp.Status != "OK" {
		return "", errors.New("upload failed: " + uploadResp.Message)
	}

	return uploadResp.Data, nil
}
