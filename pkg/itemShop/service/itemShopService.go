package service

import (
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	_playerCoinModel "github.com/Bannawat101/project-shop-api/pkg/playerCoin/model"
)

type ItemShopService interface {
	Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error)   //ดึงรายการสินค้า
	Buying(buyingReq *_itemShopModel.BuyingReq) (*_playerCoinModel.PlayerCoin, error)    //ซื้อสินค้า
	Selling(sellingReq *_itemShopModel.SellingReq) (*_playerCoinModel.PlayerCoin, error) //ขายสินค้า
}
