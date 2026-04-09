package service

import (
	"auth-project/internal/models"
	"auth-project/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwtpkg "auth-project/pkg/jwt"
)

type AuthService struct {
	repo        *repository.UserRepository
	refreshRepo *repository.RefreshRepository

	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	repo *repository.UserRepository,
	refreshRepo *repository.RefreshRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {

	return &AuthService{
		repo:        repo,
		refreshRepo: refreshRepo,
		jwtSecret:   jwtSecret,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
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

	access, err := jwtpkg.GenerateAccess(user.ID, s.jwtSecret, s.accessTTL)
	if err != nil {
		return "", "", err
	}

	refresh, err := jwtpkg.GenerateRefresh(user.ID, s.jwtSecret, s.refreshTTL)
	if err != nil {
		return "", "", err
	}

	err = s.refreshRepo.Save(
		user.ID,
		refresh,
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

	session, err := s.refreshRepo.GetSession(oldToken)
	if err != nil {
		return "", "", err
	}

	err = s.refreshRepo.Delete(oldToken)
	if err != nil {
		return "", "", err
	}

	newAccess, err := jwtpkg.GenerateAccess(userID, s.jwtSecret, s.accessTTL)
	newRefresh, err := jwtpkg.GenerateRefresh(userID, s.jwtSecret, s.refreshTTL)

	err = s.refreshRepo.Save(
		userID,
		newRefresh,
		session.DeviceID,
		session.UserAgent,
		session.IP,
		time.Now().Add(s.refreshTTL),
	)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *AuthService) DeleteRefreshToken(token string) error {
	return s.refreshRepo.Delete(token)
}
