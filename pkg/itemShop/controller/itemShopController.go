package controller

import "github.com/labstack/echo/v4"

type ItemShopController interface {
	Listing(pctx echo.Context) error //มันจะพิเศษตรงที่ต้องมี2 paramitor error.Context and error เท่านั้ย
}
