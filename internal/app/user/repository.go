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
	InitializeRoleTable() error

	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	DeleteUserByEmail(ctx context.Context, userData User) error
	GetProfileDetails(ctx context.Context, id int) (*UserProfileDetails, error)
	UpdateProfile(ctx context.Context, updateUser UserProfileDetails, id int) error
	UserApproval(ctx context.Context, email string) error
	CheckUserRole(ctx context.Context, userID uint) (bool, error)
	CreateUserRoles(ctx context.Context, userRoles UserRole) error
	ListAllUser(ctx context.Context) ([]User, error)
	UnBlockUser(ctx context.Context, id int) error
	BlockUser(ctx context.Context, id int) error
}

func NewRepository(db *gorm.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) InitializeRoleTable() error {
	migrator := r.db.Migrator()

	if !migrator.HasTable(&Role{}) {
		if err := migrator.CreateTable(&Role{}); err != nil {
			return err
		}
	}

	roles := []Role{
		{ID: 1, RoleName: "user"},
		{ID: 2, RoleName: "admin"},
		{ID: 3, RoleName: "super-admin"},
	}

	for _, role := range roles {
		var existingRole Role
		result := r.db.Where("id = ?", role.ID).First(&existingRole)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				if err := r.db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return result.Error // Only return other errors
			}
		}
	}
	return nil
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
	res := r.db.Table("users").Select("id,created_at,updated_at,deleted_at,username,email,phone_number,first_name,last_name,date_of_birth,gender,is_verified").Where("id= ?", id).First(&userData)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
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

func (r *repository) ListAllUser(ctx context.Context) ([]User, error) {
	user := []User{}
	result := r.db.Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *repository) BlockUser(ctx context.Context, id int) error {
	userData := User{}
	res := r.db.First(&userData, id)
	if res.Error != nil {
		return res.Error
	}
	if !userData.IsVerified {
		return errors.New("user already blocked")
	}

	userData.IsVerified = false
	err := r.db.Save(&userData)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *repository) UnBlockUser(ctx context.Context, id int) error {
	userData := User{}
	res := r.db.First(&userData, id)
	if res.Error != nil {
		return res.Error
	}
	if userData.IsVerified {
		return errors.New("user already unblocked")
	}
	userData.IsVerified = true
	err := r.db.Save(&userData)
	if err != nil {
		return err.Error
	}
	return nil
}
