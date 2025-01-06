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
	Overlay string `jaon:"overlay"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body pr data.", "detail": err.Error()})
		return
	}

	baseImage, err := utils.Base64ToImage(req.Base)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base image.", "detail": err.Error()})
		return
	}

	var overlayImage image.Image
	if req.Overlay != "" {
		overlayImage, err = utils.Base64ToImage(req.Overlay)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid overlay image.", "detail": err.Error()})
			return
		}
	}

	finalImage, err := services.ProcessOverlayImage(baseImage, overlayImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing overlay image.", "detail": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", finalImage)
}
