package video

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	router *gin.RouterGroup
}

func NewVideoHandler(router *gin.RouterGroup) *VideoHandler {
	videoRouter := router.Group("/video")
	return &VideoHandler{
		router: videoRouter,
	}
}

func (h *VideoHandler) SetupRoutes() {
	h.router.POST("/upload", h.handleVideoUpload)
}

func (h *VideoHandler) handleVideoUpload(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No video file provided",
		})
		return
	}
	// TODO: Implement actual video upload logic

	// save video to public/videos folder
	filePath := "./public/videos/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save video file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Video upload successful",
		"filename": file.Filename,
	})
}
