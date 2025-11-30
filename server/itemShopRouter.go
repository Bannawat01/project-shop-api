package server

import (
	_itemShopController "github.com/Bannawat101/project-shop-api/pkg/itemShop/controller"
	_itemShopRepository "github.com/Bannawat101/project-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/Bannawat101/project-shop-api/pkg/itemShop/service"
)

func (s *echoServer) initItemManagingRouter() { //กำหนด route สำหรับจัดการสินค้า
	router := s.app.Group("/v1/item-shop") //สร้าง group ของ route เพื่อจัดการกับ item shop .Group จะช่วยจัดกลุ่ม route ที่เกี่ยวข้องกันให้อยู่ด้วยกัน

	itemShopRepository := _itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger) //สร้าง instance ของ repository service และ controller
	itemShopService := _itemShopService.NewItemShopServiceImpl(itemShopRepository)          //เชื่อมโยงกัน ระหว่าง repository กับ service
	itemShopController := _itemShopController.NewItemShopControllerImpl(itemShopService)    //เชื่อมโยงกัน ระหว่าง service กับ controller

	router.GET("", itemShopController.Listing) //หลังจากที่ทำ router เสร็จแล้ว ตัว application ยังไม่รู้ว่าเรามีการประกาศ route นี้อยู่
	//ดังนั้นเราจึงต้องบอก application ว่า route นี้มีอยู่จริง เราต้องไป inject ลง server อีกที
}
