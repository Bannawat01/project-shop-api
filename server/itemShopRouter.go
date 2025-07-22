package server

import (
	_itemShopController "github.com/Bannawat101/project-shop-api/pkg/itemShop/controller"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/Bannawat101/project-shop-api/pkg/itemShop/service"
)

func (s *echoServer) initItemShopRouter() {
	router := s.app.Group("/v1/item-shop")

	itemShopRepository := _itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)
	itemShopService := _itemShopService.NewItemShopRepositoryImpl(itemShopRepository)
	itemShopController := _itemShopController.NewItemShopControllerImpl(itemShopService)

	router.GET("", itemShopController.Listing) //หลังจากที่ทำ router เสร็จแล้ว ตัว application ยังไม่รู้ว่าเรามีการประกาศ route นี้อยู่
	//ดังนั้นเราจึงต้องบอก application ว่า route นี้มีอยู่จริง เราต้องไป inject ลง server อีกที
}
