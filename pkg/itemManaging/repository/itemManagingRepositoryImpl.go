package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemManagingException "github.com/Bannawat101/project-shop-api/pkg/itemManaging/exceptions"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ItemManagingRepositoryImpl struct {
	db     *gorm.DB
	logger echo.Logger
}

func NewItemManagingRepositoryImpl(db *gorm.DB, logger echo.Logger) ItemManagingRepository {
	return &ItemManagingRepositoryImpl{db, logger}
}

func (r *ItemManagingRepositoryImpl) Creating(itemEntity *entities.Item) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Create(itemEntity).Scan(item).Error; err != nil {
		r.logger.Errorf("Error creating item: %s", err.Error())
		return nil, &_itemManagingException.ItemCreating{}
	}
	return item, nil
}
