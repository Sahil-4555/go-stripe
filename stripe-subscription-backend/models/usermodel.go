package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	Id             uint           `gorm:"primaryKey;type:int(11);" json:"id"` // Primary Key
	Name           string         `gorm:"type:varchar(50);" json:"name"`
	Email          string         `gorm:"type:varchar(30);" json:"email"`
	Address        string         `gorm:"type:text;" json:"address"`
	City           string         `gorm:"type:varchar(255);" json:"city"`
	State          string         `gorm:"type:varchar(255);" json:"state"`
	PostalCode     string         `gorm:"type:varchar(16);" json:"postal_code"`
	Password       string         `gorm:"type:varchar(255);" json:"password"`
	Country        string         `gorm:"type:varchar(255);" json:"country"`
	StripeId       string         `gorm:"type:varchar(100);" json:"stripe_id"`
	PaymentDefault bool           `gorm:"type:bool;" json:"payment_default"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
