package database

import (
	"fiapx-api/internal/entity"
	"fiapx-api/internal/repository"
	"gorm.io/gorm"
)

type VideoRepositoryGorm struct {
	DB *gorm.DB
}

var _ repository.VideoRepository = (*VideoRepositoryGorm)(nil)

func NewVideoRepository(db *gorm.DB) *VideoRepositoryGorm {
	return &VideoRepositoryGorm{DB: db}
}

func (r *VideoRepositoryGorm) Create(video *entity.Video) error {
	return r.DB.Create(video).Error
}

func (r *VideoRepositoryGorm) FindByID(id string) (*entity.Video, error) {
	var video entity.Video
	err := r.DB.First(&video, "id = ?", id).Error
	return &video, err
}

func (r *VideoRepositoryGorm) FindAllByUserID(userID string) ([]entity.Video, error) {
	var videos []entity.Video
	err := r.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&videos).Error
	return videos, err
}

func (r *VideoRepositoryGorm) Update(video *entity.Video) error {
	return r.DB.Save(video).Error
}