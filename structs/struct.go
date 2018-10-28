package structs

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Config struct
type Config struct {
	Auth        AuthConfig
	GMail       GMailConfig
	Environment EnvironmentConfig
}

// AuthConfig struct
type AuthConfig struct {
	ActivationSalt            string
	ExpirationOfActivationKey int64
}

// GMailConfig struct
type GMailConfig struct {
	SenderAddress  string
	SenderPassword string
}

// EnvironmentConfig struct
type EnvironmentConfig struct {
	Host string
}

// Model struct
type Model struct {
	ID        uint64 `json:"id" form:"id" query:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// User struct
type (
	User struct {
		gorm.Model
		Email                     string `json:"email" form:"email" query:"email" valid:"required,email"`
		Password                  string `json:"-" form:"password" valid:"required,alphanum"`
		Activated                 uint64 `json:"activated"`
		ActivationKey             string `json:"-" query:"activation_key"`
		ExpirationOfActivationKey string `json:"-"`
	}
)
