// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"arctron.cn/arctron/arcplus/gateway/conf"
	pb "arctron.cn/arctron/arcplus/pb/member"
)

type userInfo struct {
	Name   string
	Mobile string
	Email  string
}

func MemberissueJwt(c echo.Context, m *userInfo) error {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	name := m.Name
	mobile := m.Mobile
	email := m.Email

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["mobile"] = mobile
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.SignKey))
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{
			"token": "12",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

//注册会员
func Regist(c echo.Context) error {

	m := new(userInfo)
	if err := c.Bind(&m); err != nil {
		return err
	}
	name := m.Name
	mobile := m.Mobile
	email := m.Email

	cn := pb.NewMemberClient(Connections["member"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = appendRequestID(ctx, c)
	_, err := cn.Regist(ctx, &pb.RegistRequest{Name: name, Mobile: mobile, Email: email})
	if err != nil {
		Logger.Error("call micro service 'member' error: ", err)
		return errors.New("internal service error")
	}

	return MemberissueJwt(c, m)

}

//校验token
func AuthToken(c echo.Context, mobile2 string) string {

	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	mobile1 := claims["mobile"].(string)

	if mobile1 == mobile2 {
		return "true"
	} else {
		return "false"
	}
}
