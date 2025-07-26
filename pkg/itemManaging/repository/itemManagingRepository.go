package repository

import "github.com/Bannawat101/project-shop-api/entities"

type ItemManagingRepository interface {
	Creating(itemEntity *entities.Item) (*entities.Item, error)
}
