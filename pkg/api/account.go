package api

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Account model
type Account struct {
	ID        uint      `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `json:"name,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	UserID    uint      `json:"user_id,omitempty"`
	User      User      `gorm:"foreignkey:UserID"`
}

type accountSvcer interface {
	List(userID uint) ([]Account, error)
}

func newAccountService(db *gorm.DB) *accountService {
	return &accountService{db}
}

type accountService struct {
	db *gorm.DB
}

func (svc *accountService) List(userID uint) ([]Account, error) {
	var accounts []Account
	err := svc.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}
