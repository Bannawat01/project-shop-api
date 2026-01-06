package repository

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_playerCoinModel "github.com/Bannawat101/project-shop-api/pkg/playerCoin/model"
	"gorm.io/gorm"
)

type PlayerCoinRepository interface {
	CoinAdding(playerCoinEntity *entities.PlayerCoin, tx *gorm.DB) (*entities.PlayerCoin, error)
	Showing(playerID string) (*_playerCoinModel.PlayerCoinShowing, error)
}
