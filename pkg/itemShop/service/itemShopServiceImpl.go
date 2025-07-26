package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
)

type itemShopServiceImpl struct {
	ItemShopRepository _itemShopRepository.ItemShopRepository
}

func NewItemShopServiceImpl(itemShopRepository _itemShopRepository.ItemShopRepository) ItemShopService {
	return &itemShopServiceImpl{itemShopRepository}
}

func (s *itemShopServiceImpl) Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) {
	itemList, err := s.ItemShopRepository.Listing(itemFilter)
	if err != nil {
		return nil, err
	}

	itemCounting, err := s.ItemShopRepository.Counting(itemFilter)
	if err != nil {
		return nil, err
	}

	size := itemFilter.Size
	page := itemFilter.Page

	totalPage := s.totalPageCalculation(itemCounting, size)
	result := s.toItemResultResponse(itemList, page, totalPage) //go to controller next

	return result, nil

}

func (s *itemShopServiceImpl) totalPageCalculation(totalItem int64, size int64) int64 {
	if size == 0 {
		return 1 // or handle as needed
	}
	totalPage := totalItem / size

	if totalItem%size != 0 {
		totalPage++
	}
	return totalPage
}

func (s *itemShopServiceImpl) toItemResultResponse(itemEntityList []*entities.Item, page, totalPage int64) *_itemShopModel.ItemResult {
	itemModelList := make([]*_itemShopModel.Item, 0, len(itemEntityList)) //ได้มาเป็น entity เราจึงต้องแปลงกับเป็น model ด้วยการ loop
	for _, item := range itemEntityList {
		itemModelList = append(itemModelList, item.ToItemModel()) //แปลงเป็น model ด้วยการเรียกใช้ function ToItemModel
	}

	return &_itemShopModel.ItemResult{
		Items: itemModelList,
		Paginate: _itemShopModel.PaginateResult{
			Page:      page,
			TotalPage: totalPage,
		},
	}
}
