package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
)

type ItemShopRepository interface { //กำหนด interface ของ repository
	Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) //รายการสินค้า
	Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error)           //นับจำนวนสินค้า
	FindByID(itemID uint64) (*entities.Item, error)                          //ค้นหาสินค้าตาม ID
}
