package controller

import (
	"context"
	"net/http"

	"github.com/Bannawat101/project-shop-api/pkg/custom"
	_oauth2Exception "github.com/Bannawat101/project-shop-api/pkg/oauth2/exception"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (c *googleOAuth2Controller) PlayerAuthorized(pctx echo.Context, next echo.HandlerFunc) error {
	ctx := context.Background()

	TokenSource, err := c.getTokenSource(pctx)
	if err != nil {
		return custom.CustomError(pctx, http.StatusUnauthorized, err) //สาเหตที่ตรงนี้เป็น customError เพราะว่าเป็น Handler ที่ถูกเรียกใช้จาก Middleware
	}

	if !TokenSource.Valid() {
		TokenSource, err = c.playerTokenRefreshing(pctx, TokenSource)
		if err != nil {
			return custom.CustomError(pctx, http.StatusUnauthorized, err)
		}
	}

	client := playerGoogleoAuth2.Client(ctx, TokenSource)

	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	if !c.oauth2Service.IsThisGuyIsReallyPlayer(userInfo.ID) {
		return custom.CustomError(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	pctx.Set("playerID", userInfo.ID)

	return next(pctx)
}

func (c *googleOAuth2Controller) AdminAuthorized(pctx echo.Context, next echo.HandlerFunc) error {
	ctx := context.Background()

	TokenSource, err := c.getTokenSource(pctx)
	if err != nil {
		return custom.CustomError(pctx, http.StatusUnauthorized, err) //สาเหตที่ตรงนี้เป็น customError เพราะว่าเป็น Handler ที่ถูกเรียกใช้จาก Middleware
	}

	if !TokenSource.Valid() {
		TokenSource, err = c.adminTokenRefreshing(pctx, TokenSource)
		if err != nil {
			return custom.CustomError(pctx, http.StatusUnauthorized, err)
		}
	}

	client := adminGoogleoAuth2.Client(ctx, TokenSource)

	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	if !c.oauth2Service.IsThisGuyIsReallyAdmin(userInfo.ID) {
		return custom.CustomError(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	pctx.Set("adminID", userInfo.ID)

	return next(pctx)
}

func (c *googleOAuth2Controller) adminTokenRefreshing(pctx echo.Context, token *oauth2.Token) (*oauth2.Token, error) {
	ctx := context.Background()

	updateToken, err := adminGoogleoAuth2.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, updateToken.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, updateToken.RefreshToken)

	return updateToken, nil
}

func (c *googleOAuth2Controller) playerTokenRefreshing(pctx echo.Context, token *oauth2.Token) (*oauth2.Token, error) {
	ctx := context.Background()

	updateToken, err := playerGoogleoAuth2.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, updateToken.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, updateToken.RefreshToken)

	return updateToken, nil
}

func (c *googleOAuth2Controller) getTokenSource(pctx echo.Context) (*oauth2.Token, error) {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	refreshToken, err := pctx.Cookie(refreshTokenCookieName)
	if err != nil {
		return nil, &_oauth2Exception.Unauthorized{}
	}

	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}
