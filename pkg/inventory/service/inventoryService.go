package service

import _inventoryModel "github.com/Bannawat101/project-shop-api/pkg/inventory/model"

type InventoryService interface {
	Listing(playerID string) ([]*_inventoryModel.Inventory, error)
}
