package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Bannawat101/project-shop-api/config"
	"github.com/Bannawat101/project-shop-api/databases"
	_adminRepository "github.com/Bannawat101/project-shop-api/pkg/admin/repository"
	_oauth2Controller "github.com/Bannawat101/project-shop-api/pkg/oauth2/controller"
	_oauth2Service "github.com/Bannawat101/project-shop-api/pkg/oauth2/service"
	_playerRepository "github.com/Bannawat101/project-shop-api/pkg/player/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type echoServer struct {
	app  *echo.Echo
	db   databases.Database
	conf *config.Config
}

var (
	once   sync.Once
	server *echoServer
)

func NewEchoServer(conf *config.Config, db databases.Database) *echoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
		server = &echoServer{
			app:  echoApp,
			db:   db,
			conf: conf,
		}
	})

	return server
}

func (s *echoServer) Start() {
	corsMiddleware := getCORSMMiddleware(s.conf.Server.AllowOrigins)        //กำหนดค่า CORS middleware
	bodyLimitMiddleware := getBoddyLimitMiddleware(s.conf.Server.BodyLimit) //กำหนดค่า Body Limit middleware
	timeOutMiddleware := getTimeoutMiddleware(s.conf.Server.Timeout)        //กำหนดค่า Timeout middleware

	s.app.Use(middleware.Recover()) //ใช้ middleware เพื่อจัดการกับ panic ที่เกิดขึ้นในแอปพลิเคชัน
	s.app.Use(middleware.Logger())  //ใช้ middleware เพื่อบันทึก log ของคำขอที่เข้ามา
	s.app.Use(corsMiddleware)       //ใช้ CORS middleware
	s.app.Use(bodyLimitMiddleware)  //ใช้ Body Limit middleware
	s.app.Use(timeOutMiddleware)    //ใช้ Timeout middleware

	authorizingMiddleware := s.getAuthorizingMiddleware()

	s.app.GET("/v1/health", s.healthCheck)
	// s.app.GET("/v1/panic", func(c echo.Context) error {
	// 	panic("Panic")
	// })

	// Initialize routes
	s.initOAuth2Router()
	s.initItemShopRouter()
	s.initItemManagingRouter(authorizingMiddleware)
	s.initPlayerCoinRouter(authorizingMiddleware)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM) //กระบวนการเพื่อที่จะ shutdown er จำเป็นต้องมีสัญญาณ 3 ตัวนี้
	go s.gracefulShutdown(quitCh)

	s.httpListening()
}

func (s *echoServer) httpListening() {
	url := fmt.Sprintf(":%d", s.conf.Server.Port) //กำหนดพอร์ตที่เซิร์ฟเวอร์จะฟังคำขอ
	s.app.Logger.Infof("Starting server on %s", url)

	if err := s.app.Start(url); err != nil && err != http.ErrServerClosed { //ตรวจสอบข้อผิดพลาดที่อาจเกิดขึ้นระหว่างการเริ่มต้นเซิร์ฟเวอร์ HTTP
		s.app.Logger.Errorf("Server failed to start: %v", err)
		s.app.Logger.Fatal("Shutting down the server")
	}
}
func (s *echoServer) gracefulShutdown(quitCh chan os.Signal) {
	ctx := context.Background()

	<-quitCh
	s.app.Logger.Info("Shutting down the server...")

	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatal("Error shutting down the server:", err)
	}

}

func (s *echoServer) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// func getLoggerMiddleware() echo.MiddlewareFunc {
// 	return middleware.Logger()
// }

func getTimeoutMiddleware(timeOut time.Duration) echo.MiddlewareFunc { //กำหนดเวลา timeout ของ request
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{ //กำหนดค่า timeout
		Skipper:      middleware.DefaultSkipper, //กำหนดเงื่อนไขการข้าม middleware (ถ้าไม่ต้องการข้ามก็ใช้ DefaultSkipper)
		ErrorMessage: "Request timeout",         //ข้อความแสดงเมื่อเกิด timeout
		Timeout:      timeOut * time.Second,
	})
}

func getCORSMMiddleware(allawOrigin []string) echo.MiddlewareFunc { //มีไว้กัน client ที่ไม่ใช่ origin ของเราเข้ามาใช้ API
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allawOrigin,                                                            //กำหนดว่าอนุญาตให้ client จาก origin ไหนเข้ามาใช้ API ได้บ้าง
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},       //กำหนด method ที่อนุญาตให้ใช้ได้
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}, //กำหนด header ที่อนุญาตให้ใช้ได้
	})
}

func getBoddyLimitMiddleware(bodyLimt string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimt)
}

func (s *echoServer) getAuthorizingMiddleware() *authorizingMiddleware {
	playerRepository := _playerRepository.NewPlayerRepositoryImpl(s.db, s.app.Logger)
	adminRepository := _adminRepository.NewAdminRepositoryImpl(s.db, s.app.Logger)

	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(
		playerRepository,
		adminRepository,
	)

	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(
		oauth2Service,
		s.conf.OAuth2,
		s.app.Logger,
	)

	return &authorizingMiddleware{
		oauth2Controller: oauth2Controller,
		oauth2Conf:       s.conf.OAuth2,
		logger:           s.app.Logger,
	}
}
