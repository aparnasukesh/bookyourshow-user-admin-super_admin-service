package superadmin

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetSuperAdminByEmail(ctx context.Context, email string) (*User, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetSuperAdminByEmail(ctx context.Context, email string) (*User, error) {
	return &User{Email: SUPER_ADMIN, Password: SUPER_ADMIN_PASSWORD}, nil
}
