package superadmin

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetSuperAdminByEmail(ctx context.Context, email string) (*Admin, error)
	ListAdiminRequests(ctx context.Context) ([]AdminStatus, error)
	AdminApproval(ctx context.Context, email string, isVerified bool) error
	GetAdminByEmail(ctx context.Context, email string) (*Admin, error)
	CreateAdminRoles(ctx context.Context, userRoles AdminRole) error
	DeleteAdminByEmail(ctx context.Context, adminData Admin) error
	UpdateIsVerified(ctx context.Context, email string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetSuperAdminByEmail(ctx context.Context, email string) (*Admin, error) {
	return &Admin{Email: SUPER_ADMIN, Password: SUPER_ADMIN_PASSWORD}, nil
}

func (r *repository) ListAdiminRequests(ctx context.Context) ([]AdminStatus, error) {
	adminStatus := []AdminStatus{}

	if err := r.db.Where("status = ?", "pending").Find(&adminStatus).Error; err != nil {
		return nil, err
	}
	return adminStatus, nil
}

func (r *repository) AdminApproval(ctx context.Context, email string, isVerified bool) error {
	adminStatus := &AdminStatus{}

	res := r.db.Where("email=? AND status=?", email, "pending").First(&adminStatus)
	if res.Error != nil {
		return res.Error
	}
	if isVerified {
		adminStatus.Status = "approved"
	} else {
		adminStatus.Status = "rejected"
	}

	if err := r.db.Save(&adminStatus).Error; err != nil {
		return err
	}
	return nil
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

func (r *repository) CreateAdminRoles(ctx context.Context, adminRole AdminRole) error {
	if err := r.db.Create(&adminRole).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteAdminByEmail(ctx context.Context, adminData Admin) error {
	result := r.db.Where("email = ? ", adminData.Email).Delete(adminData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) UpdateIsVerified(ctx context.Context, email string) error {
	admin := &Admin{}
	if err := r.db.Where("email=?", email).First(&admin).Error; err != nil {
		return err
	}
	admin.IsVerified = true
	if err := r.db.Save(&admin).Error; err != nil {
		return err
	}
	return nil
}
