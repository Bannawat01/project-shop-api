package repository

import (
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_itemSopException "github.com/Bannawat101/project-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	"github.com/labstack/echo/v4"
)

type ItemShopRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemShopRepositoryImpl(db databases.Database, logger echo.Logger) ItemShopRepository {
	return &ItemShopRepositoryImpl{db, logger}
}

func (r *ItemShopRepositoryImpl) Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	itemList := make([]*entities.Item, 0)

	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false) //select * from items
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%") // Filter by name if provided
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%") //
	}

	// (page -1 * size) = offset
	offset := int((itemFilter.Paginate.Page - 1) * itemFilter.Paginate.Size)
	limit := int(itemFilter.Size)

	if err := query.Offset(offset).Limit(limit).Find(&itemList).Order("id asc").Error; err != nil {
		r.logger.Errorf("Error getting item list: %s", err)
		return nil, &_itemSopException.Itemisting{} // Return custom error type for item listing
	}

	return itemList, nil //Service is next
}

func (r *ItemShopRepositoryImpl) Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error) {
	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false) //select * from items
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%") // Filter by name if provided
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%") //
	}

	// count := new(int64)
	var count int64 //ถ้าใช้อันนี้เวลามี null จะไม่เกิด panic
	if err := query.Count(&count).Error; err != nil {
		r.logger.Errorf("Counting item failed: %s", err)
		return -1, &_itemSopException.ItemCounting{} // Return custom error type for item listing
	}

	return count, nil //Service is next
}

func (r *ItemShopRepositoryImpl) FindByID(itemID uint64) (*entities.Item, error) {
	item := new(entities.Item)
	if err := r.db.Connect().First(item, itemID).Error; err != nil {
		r.logger.Errorf("Item with ID %d not found: %s", itemID, err)
		return nil, &_itemSopException.ItemNotFound{}
	}
	return item, nil

}
