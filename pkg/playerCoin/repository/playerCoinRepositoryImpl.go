package repository

import (
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_playerException "github.com/Bannawat101/project-shop-api/pkg/playerCoin/excaption"
	_playerCoinModel "github.com/Bannawat101/project-shop-api/pkg/playerCoin/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func (r *PlayerCoinRepositoryImpl) CoinAdding(playerCoinEntity *entities.PlayerCoin, tx *gorm.DB) (*entities.PlayerCoin, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	playerCoin := new(entities.PlayerCoin)
	if err := conn.Create(playerCoinEntity).Scan(playerCoin).Error; err != nil {
		r.logger.Errorf("Error adding coins to player account: %s", err.Error())
		return nil, &_playerException.CoinAdding{}
	}
	return playerCoin, nil
}

func (r *PlayerCoinRepositoryImpl) Showing(playerID string) (*_playerCoinModel.PlayerCoinShowing, error) {
	playerCoinShowing := new(_playerCoinModel.PlayerCoinShowing)

	if err := r.db.Connect().Model(&entities.PlayerCoin{}).Where(
		"player_id = ?", playerID,
	).Select(
		"player_id, sum(amount) as coin",
	).Group(
		"player_id",
	).Scan(playerCoinShowing).Error; err != nil {
		r.logger.Errorf("Error showing player coins: %s", err.Error())
		return nil, &_playerException.PlayerCoinShowing{}
	}

	return playerCoinShowing, nil
}
