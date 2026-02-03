package handler

import (
	"hackaton-service-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	VideoUC *usecase.VideoUseCase
}

func NewVideoHandler(videoUC *usecase.VideoUseCase) *VideoHandler {
	return &VideoHandler{VideoUC: videoUC}
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	userID := c.GetString("userID")

	fileHeader, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo obrigatório"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir arquivo"})
		return
	}
	defer file.Close()

	video, err := h.VideoUC.RequestUpload(userID, fileHeader.Filename, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Upload iniciado",
		"video_id": video.ID,
		"status":   video.Status,
	})
}

func (h *VideoHandler) ListVideos(c *gin.Context) {
	userID := c.GetString("userID")
	
	videos, err := h.VideoUC.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}

func (h *VideoHandler) GetDownloadLink(c *gin.Context) {
	userID := c.GetString("userID")
	videoID := c.Param("id")

	url, err := h.VideoUC.GenerateDownloadURL(userID, videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}