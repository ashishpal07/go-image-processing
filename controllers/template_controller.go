package controllers

import (
	"image"
	"net/http"
	"overlay-image/services"
	"overlay-image/utils"

	"github.com/gin-gonic/gin"
)

type TemplateRequest struct {
	Base    string `json:"base" binding:"required"`
	Overlay string `json:"overlay"`
}

// TemplateController handles HTTP POST requests for overlay image processing.
// It expects a JSON body with a required `base` field containing a base64-encoded
// string of the base image, and an optional `overlay` field containing a base64-encoded
// string of the overlay image. It decodes these images, processes them using the
// ProcessOverlayImage service, and returns the resulting image as a PNG byte stream.
// In case of errors during JSON binding, image decoding, or image processing, it
// responds with appropriate HTTP error codes and messages.
func TemplateController(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or data.", "detail": err.Error()})
		return
	}

	baseImage, _, err := utils.FetchImageFromURL(req.Base)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base image URL.", "detail": err.Error()})
		return
	}

	var overlayImage image.Image
	if req.Overlay != "" {
		overlayImage, _, err = utils.FetchImageFromURL(req.Overlay)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid overlay image URL.", "detail": err.Error()})
			return
		}
	}

	uploadedImageURL, err := services.ProcessOverlayImage(baseImage, overlayImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process or upload the overlay image.", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
        "imageUrl": uploadedImageURL,
    })
}
