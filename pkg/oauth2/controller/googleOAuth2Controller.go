package controller

import (
	"math/rand"
	"sync"

	"github.com/Bannawat101/project-shop-api/config"
	_oauth2Service "github.com/Bannawat101/project-shop-api/pkg/oauth2/service"
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

func setGoogleOAuth2Config(oauth2Conf *config.OAuth2) {
	playerGoogleoAuth2 = &oauth2.Config{
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

	adminGoogleoAuth2 = &oauth2.Config{
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

func (c *googleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
