package repository

import (
	"chattrix/models"

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
  return u.db.Save(user).Error
}