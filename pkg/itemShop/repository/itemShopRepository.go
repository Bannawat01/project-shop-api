package repository

import "github.com/Bannawat101/project-shop-api/entities"

type ItemShopRepository interface { //เผื่อในกรณีที่ต้องการทำ mock
	Listing() ([]*entities.Item, error)
}
