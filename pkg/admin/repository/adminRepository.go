package repository

import "github.com/Bannawat101/project-shop-api/entities"

type AdminRepository interface {
	Creating(adminEntity *entities.Admin) (*entities.Admin, error)
	FindByID(adminID string) (*entities.Admin, error)
}
