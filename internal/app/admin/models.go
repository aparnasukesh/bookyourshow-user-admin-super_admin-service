package admin

import (
	"time"

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
	Otp         string `json:"otp"`
	Roles       []Role `gorm:"many2many:admin_roles;" json:"roles"`
}

type AdminProfileDetails struct {
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
	Title       string  `json:"title" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:text"`
	Duration    int     `json:"duration" gorm:"not null"`
	Genre       string  `json:"genre" gorm:"type:varchar(50)"`
	ReleaseDate string  `json:"release_date" gorm:"not null"`
	Rating      float64 `json:"rating" gorm:"type:decimal(3,1)"`
	Language    string  `json:"language" gorm:"type:varchar(100);not null"`
}

type TheaterType struct {
	ID              int    `json:"id"`
	TheaterTypeName string `json:"theater_type_name"`
}

type ScreenType struct {
	ID             int    `json:"id"`
	ScreenTypeName string `json:"theater_type_name"`
}

type SeatCategory struct {
	ID               int    `json:"id"`
	SeatCategoryName string `json:"seat_category_name"`
}
type Theater struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Place           string `json:"place"`
	City            string `json:"city"`
	District        string `json:"district"`
	State           string `json:"state"`
	OwnerID         uint   `json:"owner_id"`
	NumberOfScreens int    `json:"number_of_screens"`
	TheaterTypeID   int    `json:"theater_type_id"`
}

type TheaterScreen struct {
	ID           uint `json:"id"`
	TheaterID    int  `json:"theater_id"`
	ScreenNumber int  `json:"screen_number"`
	SeatCapacity int  `json:"seat_capacity"`
	ScreenTypeID int  `json:"screen_type_id"`
}

type Showtime struct {
	ID       uint      `json:"id"`
	MovieID  int       `json:"movie_id"`
	ScreenID int       `json:"screen_id"`
	ShowDate time.Time `json:"show_date"`
	ShowTime time.Time `json:"show_time"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=12"`
}
