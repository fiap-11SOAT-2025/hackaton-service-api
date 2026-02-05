package usecase_test

import (
	"errors"
	"hackaton-service-api/internal/entity"
	"hackaton-service-api/internal/usecase"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_Register(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		setup    func(m *MockUserRepository)
		wantErr  string
	}{
		{"Usuário já existe", "adm", "a@a.com", func(m *MockUserRepository) {
			m.On("FindByUsername", "adm").Return(&entity.User{}, nil)
		}, "usuário já existe"},
		{"Email já cadastrado", "new", "exists@a.com", func(m *MockUserRepository) {
			m.On("FindByUsername", "new").Return(nil, nil)
			m.On("FindByEmail", "exists@a.com").Return(&entity.User{}, nil)
		}, "email já cadastrado"},
		{"Sucesso", "valido", "v@v.com", func(m *MockUserRepository) {
			m.On("FindByUsername", "valido").Return(nil, nil)
			m.On("FindByEmail", "v@v.com").Return(nil, nil)
			m.On("Create", mock.Anything).Return(nil)
		}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepository)
			tt.setup(repo)
			uc := usecase.NewUserUseCase(repo, nil)
			err := uc.Register(tt.username, tt.email, "123456")
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserUseCase_Login(t *testing.T) {
	user, _ := entity.NewUser("test", "t@t.com", "secret")
	
	t.Run("Senha incorreta", func(t *testing.T) {
		repo := new(MockUserRepository)
		uc := usecase.NewUserUseCase(repo, nil)
		repo.On("FindByUsername", "test").Return(user, nil)
		_, _, err := uc.Login("test", "errada")
		assert.EqualError(t, err, "credenciais inválidas")
	})

	t.Run("Usuário inexistente", func(t *testing.T) {
        repo := new(MockUserRepository)
        uc := usecase.NewUserUseCase(repo, nil)

        repo.On("FindByUsername", "fantasma").Return(nil, errors.New("not found"))

        token, username, err := uc.Login("fantasma", "123")

        assert.Error(t, err)
        assert.Equal(t, "credenciais inválidas", err.Error())
        assert.Empty(t, token)
        assert.Empty(t, username)
    })
	
	t.Run("Erro no token", func(t *testing.T) {
		repo, tokenGen := new(MockUserRepository), new(MockTokenGenerator)
		uc := usecase.NewUserUseCase(repo, tokenGen)
		repo.On("FindByUsername", "test").Return(user, nil)
		tokenGen.On("GenerateToken", user.ID).Return("", errors.New("jwt error"))
		_, _, err := uc.Login("test", "secret")
		assert.Error(t, err)
	})

	t.Run("Sucesso login", func(t *testing.T) {
        repo := new(MockUserRepository)
        tokenGen := new(MockTokenGenerator)
        uc := usecase.NewUserUseCase(repo, tokenGen)

        user, _ := entity.NewUser("test", "test@teste.com", "senha123")
        
        repo.On("FindByUsername", "test").Return(user, nil)
        tokenGen.On("GenerateToken", user.ID).Return("token-valido", nil)

        token, username, err := uc.Login("test", "senha123")

        assert.NoError(t, err)
        assert.Equal(t, "token-valido", token)
        assert.Equal(t, "test", username)
    })
}