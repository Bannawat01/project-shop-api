package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Bannawat101/project-shop-api/config"
	_adminModel "github.com/Bannawat101/project-shop-api/pkg/admin/model"
	"github.com/Bannawat101/project-shop-api/pkg/custom"
	_oauth2Exception "github.com/Bannawat101/project-shop-api/pkg/oauth2/exception"
	_oauth2Model "github.com/Bannawat101/project-shop-api/pkg/oauth2/model"
	_oauth2Service "github.com/Bannawat101/project-shop-api/pkg/oauth2/service"
	_playerModel "github.com/Bannawat101/project-shop-api/pkg/player/model"
	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type googleOAuth2Controller struct {
	oauth2Service _oauth2Service.OAuth2Service
	oauth2Conf    *config.OAuth2
	logger        echo.Logger
}

var (
	playerGoogleoAuth2 *oauth2.Config //นี่คือการประกาศตัวแปรแบบ singleton สำหรับการตั้งค่า OAuth2 ของผู้เล่น
	adminGoogleoAuth2  *oauth2.Config //นี่คือการประกาศตัวแปรแบบ singleton สำหรับการตั้งค่า OAuth2 ของแอดมิน
	once               sync.Once      //ตัวแปรนี้ใช้เพื่อให้แน่ใจว่าการตั้งค่า OAuth2 จะถูกสร้างขึ้นเพียงครั้งเดียว

	accessTokenCookieName  = "act"   //ชื่อตัวแปรสำหรับคุกกี้ที่เก็บ access token
	refreshTokenCookieName = "rtc"   //ชื่อตัวแปรสำหรับคุกกี้ที่เก็บ refresh token
	stateCookieName        = "state" //ชื่อตัวแปรสำหรับคุกกี้ที่เก็บ state

	letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") //ตัวแปรนี้เก็บชุดตัวอักษรที่ใช้ในการสร้างสตริงสุ่ม
)

func NewGoogleOAuth2Controller(oauth2Service _oauth2Service.OAuth2Service, oauth2Conf *config.OAuth2, logger echo.Logger) OAuth2Controller {

	once.Do(func() {
		setGoogleOAuth2Config(oauth2Conf)
	})

	return &googleOAuth2Controller{
		oauth2Service: oauth2Service,
		oauth2Conf:    oauth2Conf,
		logger:        logger,
	}
}

func setGoogleOAuth2Config(oauth2Conf *config.OAuth2) { //ฟังก์ชั่นนี้มีขึ้นมาเพื่อกำหนดค่าการตั้งค่า OAuth2 สำหรับผู้เล่นและแอดมิน
	playerGoogleoAuth2 = &oauth2.Config{ //นี่คือการตั้งค่า OAuth2 สำหรับผู้เล่น โดยใช้ข้อมูลจาก oauth2Conf ที่ถูกส่งเข้ามา
		ClientID:     oauth2Conf.ClientId,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.PlayerRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.Endpoints.AuthUrl,
			TokenURL:      oauth2Conf.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}

	adminGoogleoAuth2 = &oauth2.Config{ //นี่คือการตั้งค่า OAuth2 สำหรับแอดมิน โดยใช้ข้อมูลจาก oauth2Conf ที่ถูกส่งเข้ามา
		ClientID:     oauth2Conf.ClientId,
		ClientSecret: oauth2Conf.ClientSecret,
		RedirectURL:  oauth2Conf.AdminRedirectUrl,
		Scopes:       oauth2Conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Conf.Endpoints.AuthUrl,
			TokenURL:      oauth2Conf.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Conf.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}

func (c *googleOAuth2Controller) PlayerLogin(pctx echo.Context) error { //ฟังก์ชั่นนี้มีขึ้นมาเพื่อจัดการกับการเข้าสู่ระบบของผู้เล่นผ่าน Google OAuth2 โดยจะสร้าง state แบบสุ่ม เก็บไว้ในคุกกี้ และเปลี่ยนเส้นทางผู้ใช้ไปยัง URL การอนุญาตของ Google OAuth2 พร้อมกับ state นั้น
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)
	return pctx.Redirect(http.StatusFound, playerGoogleoAuth2.AuthCodeURL(state)) //redirect คือ การเปลี่ยนเส้นทางผู้ใช้ไปยัง URL อื่น โดยในที่นี้คือ URL การอนุญาตของ Google OAuth2 พร้อมกับ state ที่สร้างขึ้นแบบสุ่ม

}

func (c *googleOAuth2Controller) AdminLogin(pctx echo.Context) error {
	state := c.randomState()
	c.setCookie(pctx, stateCookieName, state)
	return pctx.Redirect(http.StatusFound, adminGoogleoAuth2.AuthCodeURL(state))

}

func (c *googleOAuth2Controller) PlayerCallback(pctx echo.Context) error {
	ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callBackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Error("Error during callback validation: ", err.Error())
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	token, err := playerGoogleoAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Error("Error exchanging code for token: ", err)
		return custom.CustomError(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := playerGoogleoAuth2.Client(ctx, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Error("Error getting user info: ", err)
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	playerCreatingReq := &_playerModel.PlayerCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.PlayerAccountCreating(playerCreatingReq); err != nil {
		c.logger.Error("Error creating player account: ", err.Error())
		return custom.CustomError(pctx, http.StatusInternalServerError, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LoginResponse{Message: "login successful"})

}

func (c *googleOAuth2Controller) AdminCallback(pctx echo.Context) error {
	ctx := context.Background()
	if err := retry.Do(func() error {
		return c.callBackValidating(pctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Error("Error during callback validation: ", err.Error())
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	token, err := adminGoogleoAuth2.Exchange(ctx, pctx.QueryParam("code"))
	if err != nil {
		c.logger.Error("Error exchanging code for token: ", err)
		return custom.CustomError(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}

	client := adminGoogleoAuth2.Client(ctx, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Error("Error getting user info: ", err)
		return custom.CustomError(pctx, http.StatusUnauthorized, err)
	}

	adminCreatingReq := &_adminModel.AdminCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	if err := c.oauth2Service.AdminAccountCreating(adminCreatingReq); err != nil {
		c.logger.Error("Error creating Admin account: ", err.Error())
		return custom.CustomError(pctx, http.StatusInternalServerError, &_oauth2Exception.OAuth2Processing{})
	}

	c.setSameSiteCookie(pctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(pctx, refreshTokenCookieName, token.RefreshToken)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LoginResponse{Message: "login successful"})
}

func (c *googleOAuth2Controller) Logout(pctx echo.Context) error {
	accessToken, err := pctx.Cookie(accessTokenCookieName)
	if err != nil {
		c.logger.Error("Error retrieving access token from cookie: ", err)
		return custom.CustomError(pctx, http.StatusUnauthorized, &_oauth2Exception.Unauthorized{})
	}
	if err := c.revokeToken(accessToken.Value); err != nil {
		c.logger.Error("Error revoking token: ", err)
		return custom.CustomError(pctx, http.StatusInternalServerError, &_oauth2Exception.OAuth2Processing{})
	}

	c.removeCookie(pctx, accessTokenCookieName)
	c.removeCookie(pctx, refreshTokenCookieName)

	return pctx.JSON(http.StatusOK, &_oauth2Model.LogoutResponse{Message: "logout successful"})
}

func (c *googleOAuth2Controller) revokeToken(assessToken string) error {
	revokeUrl := fmt.Sprintf("%s?token=%s", c.oauth2Conf.RevokeUrl, assessToken)

	resp, err := http.Post(revokeUrl, "application/x-www-form-urlencoded", nil)
	if err != nil {
		c.logger.Error("Error revoking token: ", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *googleOAuth2Controller) setCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) setSameSiteCookie(pctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) removeSameSiteCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		// Secure:   true,
	}
	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) removeCookie(pctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	pctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (c *googleOAuth2Controller) callBackValidating(pctx echo.Context) error {
	state := pctx.QueryParam("state")

	stateFromCookie, err := pctx.Cookie(stateCookieName)
	if err != nil {
		c.logger.Error("Error retrieving state from cookie: ", err)
		return &_oauth2Exception.Unauthorized{}
	}

	if state != stateFromCookie.Value {
		c.logger.Error("State mismatch: expected ", stateFromCookie.Value, " but got ", state)
		return &_oauth2Exception.Unauthorized{}
	}

	return nil
}

func (c *googleOAuth2Controller) getUserInfo(client *http.Client) (*_oauth2Model.UserInfo, error) {
	resp, err := client.Get(c.oauth2Conf.UserInfoUrl)
	if err != nil {
		c.logger.Error("Error getting user info: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	userInfoInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Error reading user info response body: ", err)
		return nil, err
	}

	userInfo := new(_oauth2Model.UserInfo)
	if err := json.Unmarshal(userInfoInBytes, userInfo); err != nil {
		c.logger.Error("Error unmarshaling user info: ", err)
		return nil, err
	}

	return userInfo, nil
}
