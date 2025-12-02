package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_playerCoinModel "github.com/Bannawat101/project-shop-api/pkg/playerCoin/model"
	_playerCoinRepository "github.com/Bannawat101/project-shop-api/pkg/playerCoin/repository"
)

type playerServideCoinImpl struct {
	PlayerCoinRepository _playerCoinRepository.PlayerCoinRepository
}

func NewPlayerCoinServiceImpl(playerCoinRepository _playerCoinRepository.PlayerCoinRepository) PlayerCoinService {
	return &playerServideCoinImpl{
		PlayerCoinRepository: playerCoinRepository,
	}
}

func (s *playerServideCoinImpl) CoinAdding(coinAddingReq *_playerCoinModel.CoinAddingReq) (*_playerCoinModel.PlayerCoin, error) {
	playerCoinEntity := &entities.PlayerCoin{
		PlayerID: coinAddingReq.PlayerID,
		Amount:   coinAddingReq.Amount,
	}

	playerCoinEntityResult, err := s.PlayerCoinRepository.CoinAdding(playerCoinEntity)
	if err != nil {
		return nil, err
	}

	return playerCoinEntityResult.ToPlayerCoinModel(), nil
}
