package server

import (
	_itemManagingController "github.com/Bannawat101/project-shop-api/pkg/itemManaging/controller"
	_itemManagingRepository "github.com/Bannawat101/project-shop-api/pkg/itemManaging/repository"
	_itemManagingService "github.com/Bannawat101/project-shop-api/pkg/itemManaging/service"
)

func (s *echoServer) initItemShopRouter() {
	router := s.app.Group("/v1/item-managing")

	itemManagingRepository := _itemManagingRepository.NewItemManagingRepositoryImpl(s.db, s.app.Logger)
	itemManagingService := _itemManagingService.NewItemManagingService(itemManagingRepository)
	itemManagingController := _itemManagingController.NewItemManagingController(itemManagingService)

	router.POST("", itemManagingController.Creating)
}
