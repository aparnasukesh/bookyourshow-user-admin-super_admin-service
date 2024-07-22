package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	DeleteUserByEmail(ctx context.Context, userData User) error
	GetProfileDetails(ctx context.Context, id int) (*UserProfileDetails, error)
	UpdateProfile(ctx context.Context, updateUser UserProfileDetails, id int) error
	UserApproval(ctx context.Context, email string) error
	CheckUserRole(ctx context.Context, userID uint) (bool, error)
	CreateUserRoles(ctx context.Context, userRoles UserRole) error
}

func NewRepository(db *gorm.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	dbData := &User{}
	result := r.db.Where("email = ?", email).First(&dbData)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return dbData, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*User, error) {
	userData := User{}
	res := r.db.Table("users").Select("id,created_at,updated_at,deleted_at,username,email,phone_number,first_name,last_name,dateofbirth,gender").Where("id= ?", id).First(&userData)
	if res.Error != nil {
		return nil, res.Error
	}
	return &userData, nil
}

func (r *repository) DeleteUserByEmail(ctx context.Context, userData User) error {
	result := r.db.Where("email = ? ", userData.Email).Delete(userData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) GetProfileDetails(ctx context.Context, id int) (*UserProfileDetails, error) {
	user := &User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &UserProfileDetails{
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Gender:      user.Gender,
	}, nil
}

func (r *repository) UpdateProfile(ctx context.Context, updateUser UserProfileDetails, id int) error {
	r.db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)

	result := r.db.Model(&User{}).Where("id = ?", id).Updates(updateUser)
	if result.Error != nil {
		return result.Error
	}
	updatedUser := User{}
	if err := r.db.Where("id=?", id).First(&updatedUser).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) UserApproval(ctx context.Context, email string) error {
	userData := User{}
	res := r.db.Where("email = ?", email).First(&userData)
	if res.Error != nil {
		return res.Error
	}
	userData.IsVerified = true
	if err := r.db.Save(&userData).Error; err != nil {
		return err
	}
	return nil
}
func (r *repository) CheckUserRole(ctx context.Context, userID uint) (bool, error) {
	var userRole UserRole
	if err := r.db.Where("user_id = ? AND role_id = ?", userID, USER_ROLE).First(&userRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.New("no user exist")
		}
		return false, err
	}
	return true, nil
}

func (r *repository) CreateUserRoles(ctx context.Context, userRoles UserRole) error {
	if err := r.db.Create(&userRoles).Error; err != nil {
		return err
	}
	return nil
}
