package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"
	_itemManagingRepository "github.com/Bannawat101/project-shop-api/pkg/itemManaging/repository"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
)

type itemManagingServiceImpl struct {
	itemManagingRepository _itemManagingRepository.ItemManagingRepository
}

func NewItemManagingService(
	itemManagingRepository _itemManagingRepository.ItemManagingRepository,
) ItemManagingService {
	return &itemManagingServiceImpl{itemManagingRepository}
}

func (s *itemManagingServiceImpl) Creating(itemCreatingReq *_itemManagingModel.ItemCreateReq) (*_itemShopModel.Item, error) {
	itemEntity := &entities.Item{
		Name:        itemCreatingReq.Name,
		Description: itemCreatingReq.Description,
		Picture:     itemCreatingReq.Picture,
		Price:       uint(itemCreatingReq.Price),
	}

	itemEntity, err := s.itemManagingRepository.Creating(itemEntity)
	if err != nil {
		return nil, err
	}
	return itemEntity.ToItemModel(), nil
}
