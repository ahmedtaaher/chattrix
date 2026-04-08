package repository

import (
	"chattrix/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func(u *UserRepository) Create(user *models.User) error {
  return u.db.Create(user).Error
}

func(u *UserRepository) GetByUsername(username string) (*models.User, error) {
  var user models.User
  err := u.db.Where("username = ?", username).First(&user).Error
  return &user, err
}

func(u *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
  var user models.User
  err := u.db.First(&user, "id = ?", id).Error
  return &user, err
}

func(u *UserRepository) Update(user *models.User) error {
  return u.db.Model(&models.User{}).Where("id = ?", user.ID).Update("nickname", user.Nickname).Error
}

func(u *UserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
  return u.db.Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error
}

func(u *UserRepository) SetOnline(userID uuid.UUID) error {
  return u.db.Model(&models.User{}).
    Where("id = ?", userID).
    Updates(map[string]interface{}{
    "is_online": true,
    "last_seen": nil,
  }).Error
}

func (u *UserRepository) SetOffline(userID uuid.UUID) error {
	now := time.Now()

	return u.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_online": false,
			"last_seen": now,
		}).Error
}