package repository

import (
	"context"

	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_adminException "github.com/Bannawat101/project-shop-api/pkg/admin/exception"
	"gorm.io/gorm/logger"
)

type adminRepositoryImpl struct {
	db     databases.Database
	logger logger.Interface
}

func NewAdminRepositoryImpl(db databases.Database, logger logger.Interface) AdminRepository {
	return &adminRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *adminRepositoryImpl) Creating(adminEntity *entities.Admin) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Create(adminEntity).Scan(admin).Error; err != nil {
		r.logger.Error(context.Background(), "Error creating Admin: %s", err.Error())
		return nil, &_adminException.AdminCreating{AdminID: adminEntity.ID}
	}
	return admin, nil
}

func (r *adminRepositoryImpl) FindByID(adminID string) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Where("id = ?", adminID).First(admin).Error; err != nil {
		r.logger.Error(context.Background(), "Error finding Admin by ID: %s", err.Error())
		return nil, &_adminException.AdminCreating{AdminID: adminID}
	}
	return admin, nil
}
