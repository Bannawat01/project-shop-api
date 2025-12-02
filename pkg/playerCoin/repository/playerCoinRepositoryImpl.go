package repository

import (
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_playerException "github.com/Bannawat101/project-shop-api/pkg/playerCoin/excaption"
	"github.com/labstack/echo/v4"
)

type PlayerCoinRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewPlayerCoinRepositoryImpl(db databases.Database, logger echo.Logger) PlayerCoinRepository {
	return &PlayerCoinRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *PlayerCoinRepositoryImpl) CoinAdding(playerCoinEntity *entities.PlayerCoin) (*entities.PlayerCoin, error) {
	playerCoin := new(entities.PlayerCoin)
	if err := r.db.Connect().Create(playerCoinEntity).Scan(playerCoin).Error; err != nil {
		r.logger.Error("Error adding coins to player account:", err.Error())
		return nil, &_playerException.CoinAdding{}
	}
	return playerCoin, nil
}
