package controller

import (
	_itemShopService "github.com/Bannawat101/project-shop-api/pkg/itemShop/service"
)

type itemShopController struct {
	itemShopService _itemShopService.ItemShopService
}

func NewItemShopController(itemShopService _itemShopService.ItemShopService) ItemShopController {
	return &itemShopController{itemShopService: itemShopService}
}
