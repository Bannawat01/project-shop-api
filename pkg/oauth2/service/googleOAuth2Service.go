package service

import (
	"github.com/Bannawat101/project-shop-api/entities"
	_adminMOdel "github.com/Bannawat101/project-shop-api/pkg/admin/model"
	_adminRepository "github.com/Bannawat101/project-shop-api/pkg/admin/repository"
	_playerModel "github.com/Bannawat101/project-shop-api/pkg/player/model"
	_playerRepository "github.com/Bannawat101/project-shop-api/pkg/player/repository"
)

type googleOAuth2Service struct {
	playerRepository _playerRepository.PlayerRepository
	adminRepository  _adminRepository.AdminRepository
}

func NewGoogleOAuth2Service(playerRepo _playerRepository.PlayerRepository, adminRepo _adminRepository.AdminRepository) OAuth2Service {
	return &googleOAuth2Service{
		playerRepository: playerRepo,
		adminRepository:  adminRepo,
	}
}

func (s *googleOAuth2Service) PlayerAccountCreating(playerCreatingReq *_playerModel.PlayerCreatingReq) error {
	if !s.IsThisGuyIsReallyPlayer(playerCreatingReq.ID) {
		playerEntity := &entities.Player{
			ID:     playerCreatingReq.ID,
			Name:   playerCreatingReq.Name,
			Email:  playerCreatingReq.Email,
			Avatar: playerCreatingReq.Avatar,
		}

		if _, err := s.playerRepository.Creating(playerEntity); err != nil {
			return err
		}
	}

	return nil
}

func (s *googleOAuth2Service) AdminAccountCreating(adminCreatingReq *_adminMOdel.AdminCreatingReq) error {
	if !s.IsThisGuyIsReallyAdmin(adminCreatingReq.ID) {
		adminEntity := &entities.Admin{
			ID:     adminCreatingReq.ID,
			Name:   adminCreatingReq.Name,
			Email:  adminCreatingReq.Email,
			Avatar: adminCreatingReq.Avatar,
		}

		if _, err := s.adminRepository.Creating(adminEntity); err != nil {
			return err
		}
	}

	return nil
}

func (s *googleOAuth2Service) IsThisGuyIsReallyPlayer(playerID string) bool {
	player, err := s.playerRepository.FindByID(playerID)
	if err != nil {
		return false
	}
	return player != nil
}

func (s *googleOAuth2Service) IsThisGuyIsReallyAdmin(adminID string) bool {
	admin, err := s.adminRepository.FindByID(adminID)

	if err != nil {
		return false
	}
	return admin != nil
}
