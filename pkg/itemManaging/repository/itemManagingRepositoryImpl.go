package repository

import (
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_itemManagingException "github.com/Bannawat101/project-shop-api/pkg/itemManaging/exception"
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"

	"github.com/labstack/echo/v4"
)

type itemManagingRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemManagingRepositoryImpl(db databases.Database, logger echo.Logger) ItemManagingRepository {
	return &itemManagingRepositoryImpl{db, logger}
}

func (r *itemManagingRepositoryImpl) Creating(itemEntity *entities.Item) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Connect().Create(itemEntity).Scan(item).Error; err != nil {
		r.logger.Errorf("Error creating item: %s", err.Error())
		return nil, &_itemManagingException.ItemCreating{}
	}
	return item, nil
}

func (r *itemManagingRepositoryImpl) Editing(itemID uint64, ItemEditingReq *_itemManagingModel.ItemEditingReq) (uint64, error) {
	if err := r.db.Connect().Model(&entities.Item{}).Where("id = ?", itemID).Updates(ItemEditingReq).Error; err != nil {
		r.logger.Errorf("Error editing item with ID %d: %s", itemID, err.Error())
		return 0, &_itemManagingException.ItemEditing{}
	}
	return itemID, nil
}

func (r *itemManagingRepositoryImpl) Archive(itemID uint64) error {
	if err := r.db.Connect().Table("items").Where("id = ?", itemID).Update("is_archive", true).Error; err != nil {
		r.logger.Errorf("Error archiving item with ID %d: %s", itemID, err.Error())
		return &_itemManagingException.ItemAchiving{ItemID: itemID}
	}
	return nil
}
