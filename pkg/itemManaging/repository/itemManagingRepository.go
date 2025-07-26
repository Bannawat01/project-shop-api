package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"
)

type ItemManagingRepository interface {
	Creating(itemEntity *entities.Item) (*entities.Item, error)
	Editing(itemID uint64, ItemEditingReq *_itemManagingModel.ItemEditingReq) (uint64, error)
	Archive(itemID uint64) error
}
