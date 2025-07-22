package service

import (
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
)

type itemShopServiceImpl struct {
	ItemShopRepository _itemShopRepository.ItemShopRepository
}

func NewItemShopRepositoryImpl(itemShopRepository _itemShopRepository.ItemShopRepository) ItemShopService {
	return &itemShopServiceImpl{itemShopRepository}
}

func (s *itemShopServiceImpl) Listing() ([]*_itemShopModel.Item, error) {
	itemList, err := s.ItemShopRepository.Listing()
	if err != nil {
		return nil, err
	}

	itemModelList := make([]*_itemShopModel.Item, len(itemList)) //ได้มาเป็น entity เราจึงต้องแปลงกับเป็น model ด้วยการ loop
	for _, item := range itemList {
		itemModelList = append(itemModelList, item.ToItemModel()) //แปลงเป็น model ด้วยการเรียกใช้ function ToItemModel
	}

	return itemModelList, nil //todo:go to controller next
}
