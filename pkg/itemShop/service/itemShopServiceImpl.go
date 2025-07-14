package service

import (
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
)

type itemShopServiceImpl struct {
	ItemShopRepository _itemShopRepository.ItemShopRepository
}

func NewItemShopRepositoryImpl(itemShopRepository _itemShopRepository.ItemShopRepository) ItemShopService {
	return &itemShopServiceImpl{itemShopRepository}
}
