// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"

	pb "arctron.cn/arctron/arcplus/pb/other"
)

// 短信验证码发送
func SendSms(c echo.Context) error {

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	mobile := m["mobile"].(string)

	cn := pb.NewOtherClient(Connections["other"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = appendRequestID(ctx, c)
	r, err := cn.SendSms(ctx, &pb.SendRequest{Mobile: mobile})
	if err != nil {
		Logger.Error("call micro service 'other' error: ", err)
		return errors.New("internal service error")
	}

	if r.Sended == "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result": "success",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"result": "failed",
		})
	}

}

// 校验短信验证码
func CheckCode(c echo.Context) error {

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	mobile := m["mobile"].(string)
	code := m["code"].(string)

	cn := pb.NewOtherClient(Connections["other"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = appendRequestID(ctx, c)
	r, err := cn.CheckCode(ctx, &pb.CheckRequest{Mobile: mobile, Code: code})
	if err != nil {
		Logger.Error("call micro service 'other' error: ", err)
		return errors.New("internal service error")
	}

	if r.Checked == "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result": "success",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"result": "failed",
		})
	}

}
