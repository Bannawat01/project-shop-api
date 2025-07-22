package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ItemShopRepositoryImpl struct {
	db     *gorm.DB
	logger echo.Logger
}

func NewItemShopRepositoryImpl(db *gorm.DB, logger echo.Logger) ItemShopRepository {
	return &ItemShopRepositoryImpl{db, logger}
}

func (r *ItemShopRepositoryImpl) Listing() ([]*entities.Item, error) {
	itemList := make([]*entities.Item, 0)
	if err := r.db.Find(&itemList).Error; err != nil {
		r.logger.Errorf("Error getting item list: %s", err)
		return nil, err
	}

	return itemList, nil //todo:Service is next
}
