package repository

import "fiapx-api/internal/entity"

type VideoRepository interface {
	Create(video *entity.Video) error
	FindByID(id string) (*entity.Video, error)
	FindAllByUserID(userID string) ([]entity.Video, error)
	Update(video *entity.Video) error
}

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}