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

// UploadVideo godoc
// @Summary Realiza o upload de um vídeo
// @Description Faz o upload de um arquivo de vídeo para processamento
// @Tags Videos
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param video formData file true "Arquivo de vídeo (.mp4, .mkv, .avi)"
// @Success 202 {object} map[string]string
// @Router /api/upload [post]
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

// ListVideos godoc
// @Summary Lista vídeos do usuário
// @Description Retorna todos os vídeos enviados pelo usuário logado
// @Tags Videos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Video
// @Router /api/videos [get]
func (h *VideoHandler) ListVideos(c *gin.Context) {
	userID := c.GetString("userID")
	
	videos, err := h.VideoUC.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}

// GetDownloadLink godoc
// @Summary Gera link para download do vídeo processado
// @Description Retorna uma URL assinada do S3 para baixar o arquivo ZIP com os frames. O vídeo deve estar com status 'DONE'.
// @Tags Videos
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do Vídeo"
// @Success 200 {object} map[string]string "link: http://s3.url..."
// @Failure 401 {object} map[string]string "Acesso negado"
// @Failure 404 {object} map[string]string "Vídeo não encontrado"
// @Failure 422 {object} map[string]string "Vídeo não está pronto"
// @Router /api/videos/{id}/download [get]
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