package repository

import (
	entities "github.com/Bannawat101/project-shop-api/entities"
	_itemShopModel "github.com/Bannawat101/project-shop-api/pkg/itemShop/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ItemShopRepositoryMock struct {
	mock.Mock //pass mock object เข้าไปเพื่อใช้ฟังก์ชันของ testify/mock
}

func (m *ItemShopRepositoryMock) BeginTransaction() *gorm.DB {
	args := m.Called() // ตรงนี้หมายความว่า เรามี paramitor อะไรบ้าง แล้วเช็คว่ามันถูก call หรือป่าว
	return args.Get(0).(*gorm.DB)
}

func (m *ItemShopRepositoryMock) RollbackTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *ItemShopRepositoryMock) CommitTransaction(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *ItemShopRepositoryMock) FindByID(itemID uint64) (*entities.Item, error) {
	args := m.Called(itemID) // ตรงนี้หมายความว่า เรามี paramitor อะไรบ้าง
	return args.Get(0).(*entities.Item), args.Error(1)
	// Get(0) หมายถึง return ตัวแรก ไล่ตาม index
	// Error(1) หมายถึง return ตัวที่สอง
}

func (m *ItemShopRepositoryMock) Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	args := m.Called(itemFilter)
	return args.Get(0).([]*entities.Item), args.Error(1)
}

func (m *ItemShopRepositoryMock) FindByIDList(itemIDs []uint64) ([]*entities.Item, error) {
	args := m.Called(itemIDs)
	return args.Get(0).([]*entities.Item), args.Error(1)
}

func (m *ItemShopRepositoryMock) Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error) {
	args := m.Called(itemFilter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ItemShopRepositoryMock) PurchaseHistoryRecording(
	purchasingEntity *entities.PurchaseHistory,
	tx *gorm.DB,
) (*entities.PurchaseHistory, error) {
	args := m.Called(purchasingEntity, tx)
	return args.Get(0).(*entities.PurchaseHistory), args.Error(1)
}
