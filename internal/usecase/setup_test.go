package usecase_test

import (
	"hackaton-service-api/internal/entity"
	"mime/multipart"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct{ mock.Mock }
func (m *MockUserRepository) Create(u *entity.User) error { return m.Called(u).Error(0) }
func (m *MockUserRepository) FindByUsername(n string) (*entity.User, error) {
	args := m.Called(n)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(e string) (*entity.User, error) {
	args := m.Called(e)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepository) FindByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockVideoRepository struct{ mock.Mock }
func (m *MockVideoRepository) Create(v *entity.Video) error { return m.Called(v).Error(0) }
func (m *MockVideoRepository) FindByID(id string) (*entity.Video, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entity.Video), args.Error(1)
}
func (m *MockVideoRepository) FindAllByUserID(id string) ([]entity.Video, error) {
	args := m.Called(id)
	return args.Get(0).([]entity.Video), args.Error(1)
}
func (m *MockVideoRepository) Update(v *entity.Video) error { return m.Called(v).Error(0) }

type MockTokenGenerator struct{ mock.Mock }
func (m *MockTokenGenerator) GenerateToken(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

type MockStorageService struct{ mock.Mock }
func (m *MockStorageService) UploadFile(f multipart.File, k string) error { return m.Called(f, k).Error(0) }
func (m *MockStorageService) GeneratePresignedURL(k string) (string, error) {
	args := m.Called(k)
	return args.String(0), args.Error(1)
}
func (m *MockStorageService) GetBucketName() string { return m.Called().String(0) }

type MockQueueService struct{ mock.Mock }
func (m *MockQueueService) SendMessage(id, email string) error { return m.Called(id, email).Error(0) }