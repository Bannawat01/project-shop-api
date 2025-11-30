package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
)

type itemShopServiceImpl struct {
	ItemShopRepository _itemShopRepository.ItemShopRepository //inject repository เข้ามาใช้ใน service
}

func NewItemShopServiceImpl(itemShopRepository _itemShopRepository.ItemShopRepository) ItemShopService { //สร้าง instance ของ service implementation
	return &itemShopServiceImpl{itemShopRepository}
}

func (s *itemShopServiceImpl) Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) { //ดึงรายการสินค้า
	itemList, err := s.ItemShopRepository.Listing(itemFilter) //เรียกใช้ method Listing จาก repository
	if err != nil {
		return nil, err
	}

	itemCounting, err := s.ItemShopRepository.Counting(itemFilter) //เรียกใช้ method Counting จาก repository
	if err != nil {
		return nil, err
	}

	size := itemFilter.Size //ขนาดของหน้าที่ส่งมาจาก client เพื่อใช้ในการคำนวณหน้าทั้งหมด
	page := itemFilter.Page //หน้าปัจจุบันที่ส่งมาจาก client เพื่อใช้ในการสร้าง response

	totalPage := s.totalPageCalculation(itemCounting, size)     //คำนวณหาจำนวนหน้าทั้งหมด
	result := s.toItemResultResponse(itemList, page, totalPage) //แปลงข้อมูลเป็นรูปแบบ response ที่ต้องการ
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
