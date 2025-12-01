package controller

import "github.com/labstack/echo/v4"

type OAuth2Controller interface {
	PlayerLogin(pctx echo.Context) error //สาเหตุที่ต้องใช้ pctx echo.Context เป็นพารามิเตอร์เนื่องจากฟังก์ชั่นนี้จะถูกเรียกใช้ในบริบทของ HTTP request ซึ่ง pctx จะให้ข้อมูลและฟังก์ชั่นที่จำเป็นสำหรับการจัดการกับ request นั้น
	AdminLogin(pctx echo.Context) error
	PlayerCallback(pctx echo.Context) error
	AdminCallback(pctx echo.Context) error
	Logout(pctx echo.Context) error
}
