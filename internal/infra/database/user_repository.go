package database

import (
	"fiapx-api/internal/entity"
	"fiapx-api/internal/repository"
	"gorm.io/gorm"
)

type UserRepositoryGorm struct {
	DB *gorm.DB
}

var _ repository.UserRepository = (*UserRepositoryGorm)(nil)

func NewUserRepository(db *gorm.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{DB: db}
}

func (r *UserRepositoryGorm) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryGorm) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryGorm) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}