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
	CreateAdmin(ctx context.Context, user Admin) error
	GetAdminByEmail(ctx context.Context, email string) (*Admin, error)
	CheckIsVerified(ctx context.Context, email string) (bool, error)
	CheckAdminExist(ctx context.Context, email string) (bool, error)
	CheckAdminStatus(ctx context.Context, email string) (string, error)
	CreateAdminStatus(ctx context.Context, req *AdminStatus) error
	CheckAdminRole(ctx context.Context, userID uint) (bool, error)
	ListAllAdmin(ctx context.Context) ([]Admin, error)
	GetAdminByID(ctx context.Context, id int) (*Admin, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAdminByEmail(ctx context.Context, email string) (*Admin, error) {
	dbData := &Admin{}
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
	dbData := &Admin{}
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

func (r *repository) CreateAdmin(ctx context.Context, admin Admin) error {
	if err := r.db.Create(&admin).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CheckAdminRole(ctx context.Context, userID uint) (bool, error) {
	var userRole AdminRole
	if err := r.db.Where("admin_id = ? AND role_id = ?", userID, ADMIN_ROLE).First(&userRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.New("no admin exist")
		}
		return false, err
	}
	return true, nil
}

func (r *repository) ListAllAdmin(ctx context.Context) ([]Admin, error) {
	var admins []Admin
	result := r.db.Find(&admins)
	if result.Error != nil {
		return nil, result.Error
	}
	return admins, nil
}

func (r *repository) GetAdminByID(ctx context.Context, id int) (*Admin, error) {
	admin := &Admin{}
	result := r.db.Where("id =?", id).First(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return admin, nil
}
