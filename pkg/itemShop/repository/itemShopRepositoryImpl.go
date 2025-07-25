package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemSopException "github.com/Bannawat101/project-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
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

func (r *ItemShopRepositoryImpl) Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	itemList := make([]*entities.Item, 0)

	query := r.db.Model(&entities.Item{}) //select * from items
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%") // Filter by name if provided
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%") //
	}

	if err := query.Find(&itemList).Error; err != nil {
		r.logger.Errorf("Error getting item list: %s", err)
		return nil, &_itemSopException.Itemisting{} // Return custom error type for item listing
	}

	return itemList, nil //Service is next
}
