package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemManagingModel "github.com/Bannawat101/project-shop-api/pkg/itemManaging/model"
	_itemManagingRepository "github.com/Bannawat101/project-shop-api/pkg/itemManaging/repository"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
)

type itemManagingServiceImpl struct {
	itemManagingRepository _itemManagingRepository.ItemManagingRepository
	itemShopRepository     _itemShopRepository.ItemShopRepository
}

func NewItemManagingService(
	itemManagingRepository _itemManagingRepository.ItemManagingRepository,
	itemShopRepository _itemShopRepository.ItemShopRepository,
) ItemManagingService {
	return &itemManagingServiceImpl{
		itemManagingRepository,
		itemShopRepository,
	}
}

func (s *itemManagingServiceImpl) Creating(itemCreatingReq *_itemManagingModel.ItemCreateReq) (*_itemShopModel.Item, error) {
	itemEntity := &entities.Item{
		Name:        itemCreatingReq.Name,
		Description: itemCreatingReq.Description,
		Picture:     itemCreatingReq.Picture,
		Price:       uint(itemCreatingReq.Price),
	}

	itemEntityResult, err := s.itemManagingRepository.Creating(itemEntity)
	if err != nil {
		return nil, err
	}
	return itemEntityResult.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) Editing(itemID uint64, itemEditingReq *_itemManagingModel.ItemEditingReq) (*_itemShopModel.Item, error) {
	_, err := s.itemManagingRepository.Editing(itemID, itemEditingReq)
	if err != nil {
		return nil, err
	}
	itemEntityResult, err := s.itemShopRepository.FindByID(itemID)
	if err != nil {
		return nil, err
	}

	return itemEntityResult.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) Archive(itemID uint64) error {
	return s.itemManagingRepository.Archive(itemID)
}
