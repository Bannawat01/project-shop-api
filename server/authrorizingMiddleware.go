package server

import (
	"github.com/Bannawat101/project-shop-api/config"
	_oauth2Controller "github.com/Bannawat101/project-shop-api/pkg/oauth2/controller"
	"github.com/labstack/echo/v4"
)

type authorizingMiddleware struct {
	oauth2Controller _oauth2Controller.OAuth2Controller
	oauth2Conf       *config.OAuth2
	logger           echo.Logger
}

func (m *authorizingMiddleware) PlayerAuthorized(pctx echo.Context, next echo.HandlerFunc) echo.HandlerFunc {
	return func(pctx echo.Context) error {
		return m.oauth2Controller.PlayerAuthorized(pctx, next)
	}
}

func (m *authorizingMiddleware) AdminAuthorized(pctx echo.Context, next echo.HandlerFunc) echo.HandlerFunc {
	return func(pctx echo.Context) error {
		return m.oauth2Controller.AdminAuthorized(pctx, next)
	}
}
