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

	playerCoinEntityResult, err := s.PlayerCoinRepository.CoinAdding(playerCoinEntity, nil)
	if err != nil {
		return nil, err
	}
	playerCoinEntityResult.PlayerID = coinAddingReq.PlayerID

	return playerCoinEntityResult.ToPlayerCoinModel(), nil
}

func (s *playerServideCoinImpl) Showing(playerID string) *_playerCoinModel.PlayerCoinShowing {
	playerCoinShowing, err := s.PlayerCoinRepository.Showing(playerID)
	if err != nil {
		return &_playerCoinModel.PlayerCoinShowing{
			PlayerID: playerID,
			Coin:     0,
		}
	}
	return playerCoinShowing
}
