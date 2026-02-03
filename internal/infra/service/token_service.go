package service

import "hackaton-service-api/internal/auth"

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(userID string) (string, error) {
	return auth.GenerateToken(userID)
}

func (s *TokenService) ValidateToken(token string) (string, error) {
	return auth.ValidateToken(token)
}