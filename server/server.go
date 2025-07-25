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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type echoServer struct {
	app  *echo.Echo
	db   *gorm.DB
	conf *config.Config
}

var (
	once   sync.Once
	server *echoServer
)

func NewEchoServer(conf *config.Config, db *gorm.DB) *echoServer {
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
	corsMiddleware := getCORSMMiddleware(s.conf.Server.AllowOrigins)
	bodyLimitMiddleware := getBoddyLimitMiddleware(s.conf.Server.BodyLimit)
	timeOutMiddleware := getTimeoutMiddleware(s.conf.Server.Timeout)

	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	s.app.Use(timeOutMiddleware)

	s.app.GET("/v1/health", s.healthCheck)
	// s.app.GET("/v1/panic", func(c echo.Context) error {
	// 	panic("Panic")
	// })

	// Initialize routes
	s.initItemShopRouter()

	quitCh := make(chan os.Signal, 1)

	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM) //กระบวนการเพื่อที่จะ shutdown server จำเป็นต้องมีสัญญาณ 3 ตัวนี้
	go s.gracefulShutdown(quitCh)

	s.httpListening()
}

func (s *echoServer) httpListening() {
	url := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Infof("Starting server on %s", url)

	if err := s.app.Start(url); err != nil && err != http.ErrServerClosed {
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
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		Timeout:      timeOut * time.Second,
	})
}

func getCORSMMiddleware(allawOrigin []string) echo.MiddlewareFunc { //มีไว้กัน client ที่ไม่ใช่ origin ของเราเข้ามาใช้ API
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allawOrigin,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},       //กำหนด method ที่อนุญาตให้ใช้ได้
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}, //กำหนด header ที่อนุญาตให้ใช้ได้
	})
}

func getBoddyLimitMiddleware(bodyLimt string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimt)
}
