package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemSopException "github.com/Bannawat101/project-shop-api/pkg/itemShop/exception"
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
		return nil, &_itemSopException.Itemisting{} // Return custom error type for item listing
	}

	return itemList, nil //todo:Service is next
}
