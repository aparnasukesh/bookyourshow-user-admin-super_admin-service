package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"not null;unique" json:"username" validate:"required,min=8,max=24"`
	Password    string `gorm:"not null" json:"password" validate:"required,min=6,max=12"`
	Email       string `gorm:"not null;unique" json:"email" validate:"email,required"`
	PhoneNumber string `gorm:"not null" json:"phone" validate:"required,len=10"`
	FirstName   string `gorm:"not null" json:"firstname" validate:"required"`
	LastName    string `gorm:"not null" json:"lastname" validate:"required"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	IsVerified  bool   `gorm:"default:false" json:"is_verified"`
	Otp         string `json:"otp"`
	Roles       []Role `gorm:"many2many:user_roles;" json:"roles"`
}

type Role struct {
	ID       uint   `gorm:"primaryKey"`
	RoleName string `gorm:"size:255;not null;unique" json:"role_name" validate:"required"`
	Users    []User `gorm:"many2many:user_roles;" json:"users"`
}
type UserRole struct {
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role   Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserProfileDetails struct {
	ID          int    `json:"id"`
	Username    string `json:"username" validate:"required,min=8,max=24"`
	Password    string `json:"password" validate:"required,min=6,max=12"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	FirstName   string `json:"firstname" validate:"required"`
	LastName    string `json:"lastname" validate:"required"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	IsVerified  bool   `json:"is_verified"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=12"`
}
