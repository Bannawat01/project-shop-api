package entities

import (
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	"time"
)

type Item struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	AdminID     *string   `gorm:"type:varchar(64);"` // Nullable field for AdminID
	Name        string    `gorm:"type:varchar(64);unique;not null;"`
	Description string    `gorm:"type:varchar(128);not null;"`
	Picture     string    `gorm:"type:varchar(256);not null;"`
	Price       uint      `gorm:"not null;"`
	IsArchive   bool      `gorm:"not null;default:false;"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime;"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime;"`
}

func (i *Item) ToItemModel() *_itemShopModel.Item { //ต้องมี function นี้เพื่อที่จะแปลงกลับไปเป็น item shop model
	return &_itemShopModel.Item{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		Picture:     i.Picture,
		Price:       i.Price,
	}
}
