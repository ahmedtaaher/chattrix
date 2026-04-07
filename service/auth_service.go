package service

import (
	"chattrix/dto"
	"chattrix/models"
	"chattrix/repository"
	"chattrix/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtService *utils.JWTService
}

func NewAuthService(repo *repository.UserRepository, jwtService *utils.JWTService) *AuthService {
	return &AuthService{repo: repo, jwtService: jwtService}
}

func (a *AuthService) Register(registerRequest dto.RegisterRequest) error {
	_, err := a.repo.GetByUsername(registerRequest.Username)
	if err == nil {
		return errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check username: %w", err)
	}

	hashedPassword, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username:     registerRequest.Username,
		Nickname:     registerRequest.Nickname,
		PasswordHash: hashedPassword,
	}

	if err := a.repo.Create(&user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (a *AuthService) Login(loginRequest dto.LoginRequest) (string, error) {
	user, err := a.repo.GetByUsername(loginRequest.Username)
	if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return "", errors.New("invalid username or password")
    }
    return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	token, err := a.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a *AuthService) GetProfile(userID uuid.UUID) (*models.User, error) {
	user, err := a.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

func(a *AuthService) UpdateProfile(userID uuid.UUID, updateRequest dto.UpdateProfileRequest) error {
  user, err := a.repo.GetByID(userID)
  if err != nil {
    return fmt.Errorf("failed to fetch user: %w", err)
  }

  if updateRequest.Nickname != "" {
    user.Nickname = updateRequest.Nickname
  }

  return a.repo.Update(user)
}

func(a *AuthService) ChangePassword(userID uuid.UUID, changePasswordRequest dto.ChangePasswordRequest) error {
  user, err := a.repo.GetByID(userID)
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return errors.New("user not found")
    }
    return fmt.Errorf("failed to fetch user: %w", err)
  }

  if !utils.CheckPasswordHash(changePasswordRequest.CurrentPassword, user.PasswordHash) {
		return errors.New("invalid password")
	}

  if utils.CheckPasswordHash(changePasswordRequest.NewPassword, user.PasswordHash) {
    return errors.New("new password cannot be the same as the current password")
  }

  hash, err := utils.HashPassword(changePasswordRequest.NewPassword)
  if err != nil {
    return fmt.Errorf("failed to hash new password: %w", err)
  }

  user.PasswordHash = hash

  if err:= a.repo.Update(user); err != nil {
    return fmt.Errorf("failed to update password: %w", err)
  }

  return nil
}

func (a *AuthService) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return a.repo.UpdateAvatar(userID, avatarURL)
}
