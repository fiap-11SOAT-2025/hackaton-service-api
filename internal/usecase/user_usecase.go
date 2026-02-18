package usecase

import (
	"errors"
	"hackaton-service-api/internal/entity"
	"hackaton-service-api/internal/repository"
)

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}

type UserUseCase struct {
	Repo  repository.UserRepository
	Token TokenGenerator
}

func NewUserUseCase(repo repository.UserRepository, token TokenGenerator) *UserUseCase {
	return &UserUseCase{
		Repo:  repo,
		Token: token,
	}
}

func (uc *UserUseCase) Register(username, email, password string) error {
	existingUser, _ := uc.Repo.FindByUsername(username)
	if existingUser != nil {
		return errors.New("usuário já existe")
	}

	existingEmail, _ := uc.Repo.FindByEmail(email)
	if existingEmail != nil {
		return errors.New("email já cadastrado")
	}

	user, err := entity.NewUser(username, email, password)
	if err != nil {
		return err
	}

	return uc.Repo.Create(user)
}

func (uc *UserUseCase) Login(username, password string) (string, string, error) {
	user, err := uc.Repo.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("credenciais inválidas")
	}

	if !user.ValidatePassword(password) {
		return "", "", errors.New("credenciais inválidas")
	}

	token, err := uc.Token.GenerateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return token, user.Username, nil
}