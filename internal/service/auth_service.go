package service

import (
	"auth-project/internal/models"
	"auth-project/internal/repository"

	"golang.org/x/crypto/bcrypt"

	jwtpkg "auth-project/pkg/jwt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(email, password string) (int, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return 0, err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	id, err := s.repo.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *AuthService) Login(email, password string) (string, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", err
	}

	token, err := jwtpkg.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
