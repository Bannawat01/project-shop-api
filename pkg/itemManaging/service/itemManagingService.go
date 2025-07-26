package service

import (
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
)

type ItemManagingService interface {
	Creating(itemCreatingReq *_itemManagingModel.ItemCreateReq) (*_itemShopModel.Item, error)
	Editing(itemID uint64, itemEditingReq *_itemManagingModel.ItemEditingReq) (*_itemShopModel.Item, error)
	Archive(itemID uint64) error
}
