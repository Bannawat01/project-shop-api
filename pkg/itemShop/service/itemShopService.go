package service

import (
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
)

type ItemShopService interface {
	Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) //_itemShopModel 	ต้อง import เข้ามาเอง
}
