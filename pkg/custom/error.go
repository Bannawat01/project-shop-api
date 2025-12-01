package custom

import "github.com/labstack/echo/v4"

type ErrorMessage struct {
	Message string `json:"message"`
}

func CustomError(pctx echo.Context, statusCode int, err error) error { //ฟังก์ชันสำหรับส่งข้อความแสดงข้อผิดพลาดแบบกำหนดเอง
	return pctx.JSON(statusCode, &ErrorMessage{Message: err.Error()}) //ส่งข้อความแสดงข้อผิดพลาดในรูปแบบ JSON พร้อมกับรหัสสถานะ HTTP ที่กำหนด
}
