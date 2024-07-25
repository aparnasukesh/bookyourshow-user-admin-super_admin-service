package sql

import (
	"fmt"
	"log"
	"sync"

	"github.com/aparnasukesh/user-admin-super_admin-svc/config"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/admin"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	mutex      sync.Mutex
	isExist    map[string]bool
)

func NewSql(config config.Config) (*gorm.DB, error) {
	if dbInstance == nil && !isExist[config.DBName] {
		mutex.Lock()
		defer mutex.Unlock()

		if dbInstance == nil && !isExist[config.DBName] {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			dbInstance = db
		}
	}
	dbInstance.AutoMigrate(&user.User{})
	dbInstance.AutoMigrate(&user.Role{})
	dbInstance.AutoMigrate(&user.UserRole{})
	dbInstance.AutoMigrate(&admin.Admin{})
	dbInstance.AutoMigrate(&admin.AdminStatus{})
	dbInstance.AutoMigrate(&admin.AdminRole{})

	userRepo := user.NewRepository(dbInstance)
	if err := userRepo.InitializeRoleTable(); err != nil {
		log.Fatal(err.Error())
	}

	return dbInstance, nil
}
