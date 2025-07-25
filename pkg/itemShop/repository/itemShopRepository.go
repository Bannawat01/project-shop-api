package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
)

type ItemShopRepository interface { //เผื่อในกรณีที่ต้องการทำ mock
	Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error)
	Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error)
	FindByID(itemID uint64) (*entities.Item, error)
}
