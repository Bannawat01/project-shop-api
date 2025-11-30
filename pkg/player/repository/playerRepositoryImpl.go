package repository

import (
	"context"

	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	_PlayerException "github.com/Bannawat101/project-shop-api/pkg/player/exception"
	"gorm.io/gorm/logger"
)

type PlayerRepositoryImpl struct {
	db     databases.Database
	logger logger.Interface
}

func NewPlayerRepositoryImpl(db databases.Database, logger logger.Interface) PlayerRepository {
	return &PlayerRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *PlayerRepositoryImpl) Creating(PlayerEntity *entities.Player) (*entities.Player, error) {
	Player := new(entities.Player)

	if err := r.db.Connect().Create(PlayerEntity).Scan(Player).Error; err != nil {
		r.logger.Error(context.Background(), "Error creating Player: %s", err.Error())
		return nil, &_PlayerException.PlayerCreating{PlayerID: PlayerEntity.ID}
	}
	return Player, nil
}

func (r *PlayerRepositoryImpl) FindByID(playerID string) (*entities.Player, error) {
	player := new(entities.Player)

	if err := r.db.Connect().Where("id = ?", playerID).First(player).Error; err != nil {
		r.logger.Error(context.Background(), "Error finding Player by ID: %s", err.Error())
		return nil, &_PlayerException.PlayerCreating{PlayerID: playerID}
	}
	return player, nil
}
