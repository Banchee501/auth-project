package service

import (
	"auth-project/internal/models"
	"auth-project/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	jwtpkg "auth-project/pkg/jwt"
)

type AuthService struct {
	repo        *repository.UserRepository
	refreshRepo *repository.RefreshRepository
}

func NewAuthService(
	repo *repository.UserRepository,
	refreshRepo *repository.RefreshRepository,
) *AuthService {
	return &AuthService{
		repo:        repo,
		refreshRepo: refreshRepo,
	}
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

func (s *AuthService) Login(email, password, deviceID, userAgent, ip string) (string, string, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", err
	}

	access, _ := jwtpkg.GenerateAccess(user.ID)
	refresh, _ := jwtpkg.GenerateRefresh(user.ID)

	deviceID := uuid.New().String()

	err = s.refreshRepo.Save(
		user.ID,
		refreshToken,
		deviceID,
		userAgent,
		ip,
		time.Now().Add(s.refreshTTL),
	)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *AuthService) Refresh(oldToken string) (string, string, error) {

	userID, err := s.refreshRepo.Find(oldToken)
	if err != nil {
		return "", "", err
	}

	s.refreshRepo.Delete(oldToken)

	newAccess, _ := jwtpkg.GenerateAccess(userID)
	newRefresh, _ := jwtpkg.GenerateRefresh(userID)

	deviceID := uuid.New().String()
	s.refreshRepo.Save(userID, newRefresh, deviceID, "", "", time.Now().Add(time.Hour*24*7))

	return newAccess, newRefresh, nil
}

func (s *AuthService) DeleteRefreshToken(token string) error {
	return s.refreshRepo.Delete(token)
}
