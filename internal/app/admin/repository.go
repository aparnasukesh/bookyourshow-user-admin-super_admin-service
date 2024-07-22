package admin

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	CreateAdmin(ctx context.Context, user User) error
	GetAdminByEmail(ctx context.Context, email string) (*User, error)
	CheckIsVerified(ctx context.Context, email string) (bool, error)
	CheckAdminExist(ctx context.Context, email string) (bool, error)
	CheckAdminStatus(ctx context.Context, email string) (string, error)
	CreateAdminStatus(ctx context.Context, req *AdminStatus) error
	CheckAdminRole(ctx context.Context, userID uint) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAdminByEmail(ctx context.Context, email string) (*User, error) {
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

func (r *repository) CheckIsVerified(ctx context.Context, email string) (bool, error) {
	dbData := &User{}
	result := r.db.Where("email = ? AND is_verified =?", email, true).First(&dbData)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repository) CheckAdminExist(ctx context.Context, email string) (bool, error) {
	dbdata := &AdminStatus{}
	result := r.db.Where("email = ? ", email).First(&dbdata)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *repository) CheckAdminStatus(ctx context.Context, email string) (string, error) {
	var status AdminStatus
	result := r.db.Where("email = ?", email).First(&status)
	if result.Error != nil {
		return "", result.Error
	}

	return status.Status, nil
}

func (r *repository) CreateAdminStatus(ctx context.Context, req *AdminStatus) error {
	if err := r.db.Create(&req).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateAdmin(ctx context.Context, user User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CheckAdminRole(ctx context.Context, userID uint) (bool, error) {
	var userRole UserRole
	if err := r.db.Where("user_id = ? AND role_id = ?", userID, ADMIN_ROLE).First(&userRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.New("no admin exist")
		}
		return false, err
	}
	return true, nil
}
