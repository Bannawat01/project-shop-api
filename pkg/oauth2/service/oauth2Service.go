package service

import (
	_adminModel "github.com/Bannawat101/project-shop-api/pkg/admin/model"
	_playerModel "github.com/Bannawat101/project-shop-api/pkg/player/model"
)

type OAuth2Service interface {
	PlayerAccountCreating(playerCreatingReq *_playerModel.PlayerCreatingReq) error
	AdminAccountCreating(adminCreatingReq *_adminModel.AdminCreatingReq) error
	IsThisGuyIsReallyPlayer(playerID string) bool
	IsThisGuyIsReallyAdmin(adminID string) bool
}
