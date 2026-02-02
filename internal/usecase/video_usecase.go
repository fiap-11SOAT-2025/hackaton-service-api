package usecase

import (
	"fiapx-api/internal/entity"
	"fiapx-api/internal/repository"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
)

type FileStorageService interface {
	UploadFile(file multipart.File, key string) error
	GeneratePresignedURL(key string) (string, error)
	GetBucketName() string
}

type QueueService interface {
	SendMessage(message string) error
}

type VideoUseCase struct {
	Repo    repository.VideoRepository
	Storage FileStorageService
	Queue   QueueService
}

func NewVideoUseCase(repo repository.VideoRepository, storage FileStorageService, queue QueueService) *VideoUseCase {
	return &VideoUseCase{
		Repo:    repo,
		Storage: storage,
		Queue:   queue,
	}
}

func (uc *VideoUseCase) RequestUpload(userID string, fileName string, file multipart.File) (*entity.Video, error) {
	ext := filepath.Ext(fileName)
	if ext != ".mp4" && ext != ".mkv" && ext != ".avi" {
		return nil, fmt.Errorf("formato não suportado")
	}

	uniqueName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileName)
	s3Key := "uploads/" + uniqueName

	video := entity.NewVideo(userID, fileName, s3Key)
	video.InputBucket = uc.Storage.GetBucketName()

	if err := uc.Repo.Create(video); err != nil {
		return nil, err
	}

	if err := uc.Storage.UploadFile(file, uniqueName); err != nil {
		return nil, err
	}

	if err := uc.Queue.SendMessage(video.ID); err != nil {
		video.Status = entity.StatusError
		video.ErrorMessage = "Falha ao enfileirar"
		uc.Repo.Update(video)
		return nil, err
	}

	return video, nil
}

func (uc *VideoUseCase) ListByUser(userID string) ([]entity.Video, error) {
	return uc.Repo.FindAllByUserID(userID)
}

func (uc *VideoUseCase) GenerateDownloadURL(userID, videoID string) (string, error) {
	video, err := uc.Repo.FindByID(videoID)
	if err != nil {
		return "", err
	}

	if video.UserID != userID {
		return "", fmt.Errorf("acesso negado")
	}

	if video.Status != entity.StatusDone {
		return "", fmt.Errorf("vídeo não está pronto")
	}

	return uc.Storage.GeneratePresignedURL(video.OutputKey)
}