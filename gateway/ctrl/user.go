// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"google.golang.org/grpc/metadata"

	"arctron.cn/arctron/arcplus/gateway/conf"
	pb "arctron.cn/arctron/arcplus/pb/user"
)

func issueJwt(c echo.Context, r *pb.AuthReply) error {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = r.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.SignKey))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func appendRequestID(ctx context.Context, c echo.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
}

// Login handlers
func Login(c echo.Context) error {
	// Use the following 2 lines for form body
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	username := m["username"].(string)
	password := m["password"].(string)

	cn := pb.NewUserClient(Connections["user"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = appendRequestID(ctx, c)
	r, err := cn.Auth(ctx, &pb.AuthRequest{Username: username, Password: password})
	if err != nil {
		Logger.Error("call micro service 'user' error: ", err)
		return errors.New("internal service error")
	}

	if r.Authed {

		return issueJwt(c, r)
	}

	return echo.ErrUnauthorized
}

// OAuth third party login
func OAuth(c echo.Context) error {

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	provider := m["provider"].(string)
	openid := m["openid"].(string)

	cn := pb.NewUserClient(Connections["user"])

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	ctx = appendRequestID(ctx, c)

	r, err := cn.OAuth(ctx, &pb.OAuthRequest{Provider: provider, Openid: openid})

	if err != nil {
		Logger.Error("call micro service 'user' error: ", err)
		return errors.New("internal service error")
	}

	if r.Authed {
		return issueJwt(c, r)
	}
	return echo.ErrUnauthorized
}

// Auth handler
func Auth(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	name := claims["jti"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
