// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"google.golang.org/grpc/metadata"

	pb "arctron.cn/arctron/arcplus/pb/qiniu"
)

// 获取七牛上传token
func MediaToken(c echo.Context) error {

	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userid := m["userid"].(string)
	mobile := m["mobile"].(string)

	//校验token
	authtoken := AuthToken(c, mobile)
	if authtoken != "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result":  "false",
			"message": "token error",
		})
	}

	cn := pb.NewQiniuClient(Connections["qiniu"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	r, err := cn.GetToken(ctx, &pb.UpTokenRequest{Userid: userid})
	if err != nil {
		Logger.Error("call micro service 'qiniu' error: ", err)
		return errors.New("internal service error")
	}

	if r.Token != "" {
		return c.JSON(http.StatusOK, map[string]string{
			"token": r.Token,
		})
	}

	return echo.ErrUnauthorized
}
