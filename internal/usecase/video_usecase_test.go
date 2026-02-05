package usecase_test

import (
	"errors"
	"hackaton-service-api/internal/entity"
	"hackaton-service-api/internal/usecase"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVideoUseCase_RequestUpload(t *testing.T) {
	t.Run("Erro: Formato de arquivo não suportado", func(t *testing.T) {
		uc := usecase.NewVideoUseCase(nil, nil, nil, nil)
		video, err := uc.RequestUpload("user1", "documento.pdf", nil)

		assert.Nil(t, video)
		assert.EqualError(t, err, "formato não suportado")
	})

	t.Run("Erro: Usuário não encontrado", func(t *testing.T) {
		userRepo := new(MockUserRepository)
		uc := usecase.NewVideoUseCase(nil, userRepo, nil, nil)

		userRepo.On("FindByID", "user_fantasma").Return(nil, errors.New("not found"))

		video, err := uc.RequestUpload("user_fantasma", "video.mp4", nil)
		assert.Nil(t, video)
		assert.Contains(t, err.Error(), "usuário não encontrado")
	})

	t.Run("Erro: Falha ao criar registro no banco", func(t *testing.T) {
		repo, userRepo, storage := new(MockVideoRepository), new(MockUserRepository), new(MockStorageService)
		uc := usecase.NewVideoUseCase(repo, userRepo, storage, nil)

		userRepo.On("FindByID", "u1").Return(&entity.User{Email: "u@u.com"}, nil)
		storage.On("GetBucketName").Return("bucket")
		repo.On("Create", mock.Anything).Return(errors.New("db error"))

		video, err := uc.RequestUpload("u1", "v.mp4", nil)
		assert.Nil(t, video)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Erro: Falha no upload para o Storage", func(t *testing.T) {
		repo, userRepo, storage := new(MockVideoRepository), new(MockUserRepository), new(MockStorageService)
		uc := usecase.NewVideoUseCase(repo, userRepo, storage, nil)

		userRepo.On("FindByID", "u1").Return(&entity.User{Email: "u@u.com"}, nil)
		storage.On("GetBucketName").Return("bucket")
		repo.On("Create", mock.Anything).Return(nil)
		storage.On("UploadFile", mock.Anything, mock.Anything).Return(errors.New("s3 error"))

		video, err := uc.RequestUpload("u1", "v.mp4", nil)
		assert.Nil(t, video)
		assert.EqualError(t, err, "s3 error")
	})

	t.Run("Erro na fila SQS", func(t *testing.T) {
		repo, userRepo, storage, queue := new(MockVideoRepository), new(MockUserRepository), new(MockStorageService), new(MockQueueService)
		uc := usecase.NewVideoUseCase(repo, userRepo, storage, queue)

		userRepo.On("FindByID", "u1").Return(&entity.User{Email: "e@e.com"}, nil)
		storage.On("GetBucketName").Return("b")
		repo.On("Create", mock.Anything).Return(nil)
		storage.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
		queue.On("SendMessage", mock.Anything, "e@e.com").Return(errors.New("sqs fail"))
		repo.On("Update", mock.Anything).Return(nil) // Cobre o handleError do UseCase

		_, err := uc.RequestUpload("u1", "v.mp4", nil)
		assert.Error(t, err)
	})

	t.Run("Sucesso: Fluxo completo de upload", func(t *testing.T) {
		repo, userRepo, storage, queue := new(MockVideoRepository), new(MockUserRepository), new(MockStorageService), new(MockQueueService)
		uc := usecase.NewVideoUseCase(repo, userRepo, storage, queue)

		userRepo.On("FindByID", "u1").Return(&entity.User{Email: "e@e.com"}, nil)
		storage.On("GetBucketName").Return("bucket")
		repo.On("Create", mock.Anything).Return(nil)
		storage.On("UploadFile", mock.Anything, mock.Anything).Return(nil)
		queue.On("SendMessage", mock.Anything, "e@e.com").Return(nil) // Agora retorna nil para sucesso

		video, err := uc.RequestUpload("u1", "video.mp4", nil)
		assert.NoError(t, err)
		assert.NotNil(t, video)
		assert.Equal(t, "video.mp4", video.FileName)
	})
}

func TestVideoUseCase_ListByUser(t *testing.T) {
	repo := new(MockVideoRepository)
	uc := usecase.NewVideoUseCase(repo, nil, nil, nil)
	repo.On("FindAllByUserID", "u1").Return([]entity.Video{{ID: "1"}}, nil)

	videos, err := uc.ListByUser("u1")
	assert.NoError(t, err)
	assert.Len(t, videos, 1)
}

func TestVideoUseCase_GenerateDownloadURL(t *testing.T) {
	t.Run("Erro: Vídeo não encontrado no banco", func(t *testing.T) {
		repo := new(MockVideoRepository)
		uc := usecase.NewVideoUseCase(repo, nil, nil, nil)

		repo.On("FindByID", "v_inexistente").Return(nil, errors.New("sql: no rows"))

		url, err := repo.FindByID("v_inexistente")
		_ = url 
		assert.Error(t, err)

		res, err := uc.GenerateDownloadURL("u1", "v_inexistente")
		assert.Empty(t, res)
		assert.Error(t, err)
	})

	t.Run("Erro: Acesso negado (Usuário incorreto)", func(t *testing.T) {
		repo := new(MockVideoRepository)
		uc := usecase.NewVideoUseCase(repo, nil, nil, nil)

		video := &entity.Video{UserID: "dono_original"}
		repo.On("FindByID", "v1").Return(video, nil)

		url, err := uc.GenerateDownloadURL("hacker", "v1")
		assert.Empty(t, url)
		assert.EqualError(t, err, "acesso negado")
	})

	t.Run("Vídeo não pronto", func(t *testing.T) {
		repo := new(MockVideoRepository)
		uc := usecase.NewVideoUseCase(repo, nil, nil, nil)
		repo.On("FindByID", "v1").Return(&entity.Video{UserID: "u1", Status: entity.StatusPending}, nil)

		_, err := uc.GenerateDownloadURL("u1", "v1")
		assert.EqualError(t, err, "vídeo não está pronto")
	})

	t.Run("Sucesso: Gera URL assinada", func(t *testing.T) {
		repo, storage := new(MockVideoRepository), new(MockStorageService)
		uc := usecase.NewVideoUseCase(repo, nil, storage, nil)

		video := &entity.Video{
			UserID:    "u1",
			Status:    entity.StatusDone,
			OutputKey: "final.zip",
		}
		repo.On("FindByID", "v1").Return(video, nil)
		storage.On("GeneratePresignedURL", "final.zip").Return("http://aws-link.com/file", nil)

		url, err := uc.GenerateDownloadURL("u1", "v1")
		assert.NoError(t, err)
		assert.Equal(t, "http://aws-link.com/file", url)
	})
}