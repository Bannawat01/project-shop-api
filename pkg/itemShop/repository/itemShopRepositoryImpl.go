package repository

import (
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_itemShopException "github.com/Bannawat101/project-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type itemRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemShopRepositoryImpl(db databases.Database, logger echo.Logger) ItemShopRepository {
	return &itemRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *itemRepositoryImpl) BeginTransaction() *gorm.DB {
	tx := r.db.Connect()
	return tx.Begin()
}

func (r *itemRepositoryImpl) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *itemRepositoryImpl) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *itemRepositoryImpl) Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false)

	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	offset := int((itemFilter.Page - 1) * itemFilter.Size)
	size := int(itemFilter.Size)

	items := make([]*entities.Item, 0)

	if err := query.Offset(offset).Limit(size).Find(&items).Order("id asc").Error; err != nil {
		r.logger.Error("Failed to find items", err.Error())
		return nil, &_itemShopException.ItemListing{}
	}

	return items, nil
}

func (r *itemRepositoryImpl) Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error) {
	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false)

	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	var count int64

	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("Counting items failed:", err.Error())
		return -1, &_itemShopException.ItemCounting{}
	}

	return count, nil
}

func (r *itemRepositoryImpl) FindByID(itemID uint64) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Connect().First(item, itemID).Error; err != nil {
		r.logger.Error("Finding item failed:", err.Error())
		return nil, &_itemShopException.ItemNotFound{ItemID: itemID}
	}

	return item, nil

}

func (r *itemRepositoryImpl) FindByIDList(itemIDs []uint64) ([]*entities.Item, error) {
	items := make([]*entities.Item, 0)

	if err := r.db.Connect().Model(&entities.Item{}).Where("id in ?", itemIDs).Find(&items).Error; err != nil {
		r.logger.Error("Finding items by ID failed:", err.Error())
		return nil, &_itemShopException.ItemListing{}
	}

	return items, nil
}

func (r *itemRepositoryImpl) PurchaseHistoryRecording(
	purchasingEntity *entities.PurchaseHistory,
	tx *gorm.DB,
) (*entities.PurchaseHistory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	insertedPurchasing := new(entities.PurchaseHistory)

	if err := conn.Create(purchasingEntity).Scan(insertedPurchasing).Error; err != nil {
		r.logger.Errorf("Purchase history recording failed: %s", err.Error())
		return nil, &_itemShopException.HistoryOfPurchaseRecording{}
	}

	return insertedPurchasing, nil
}
