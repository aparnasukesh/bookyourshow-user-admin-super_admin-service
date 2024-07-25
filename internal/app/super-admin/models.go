package superadmin

import (
	"gorm.io/gorm"
)

type Admin struct {
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
	Roles       []Role `gorm:"many2many:admin_roles;" json:"roles"`
}

type Role struct {
	gorm.Model
	RoleName string  `gorm:"size:255;not null;unique" json:"role_name" validate:"required"`
	Admins   []Admin `gorm:"many2many:admin_roles;" json:"admins"`
}

type AdminRole struct {
	AdminID uint  `json:"admin_id"`
	RoleID  uint  `json:"role_id"`
	Admin   Admin `gorm:"foreignKey:AdminID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role    Role  `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AdminStatus struct {
	gorm.Model
	Status string `gorm:"not null" json:"status"`
	Email  string `gorm:"not null;unique" json:"email" validate:"email,required"`
}

type SuperAdminProfileDetails struct {
	Username    string `json:"username" validate:"required,min=8,max=24"`
	Password    string `json:"password" validate:"required,min=6,max=12"`
	PhoneNumber string `json:"phone" validate:"required,len=10"`
	Email       string `json:"email" validate:"email,required"`
	FirstName   string `json:"firstname" validate:"required"`
	LastName    string `json:"lastname" validate:"required"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

type AdminRequests struct {
	Email string `json:"email" validate:"email,required"`
}

type Movie struct {
	Title       string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text"`
	Duration    int     `gorm:"not null"`
	Genre       string  `gorm:"type:varchar(50)"`
	ReleaseDate string  `gorm:"not null"`
	Rating      float64 `gorm:"type:decimal(3,1)"`
	Language    string  `gorm:"type:varchar(100);not null"`
}
