// Package ctrl for controllers used in gateway
package ctrl

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"google.golang.org/grpc/metadata"

	pb "arctron.cn/arctron/arcplus/pb/expert"
)

type ExpertInfo struct {
	Userid      string
	Pic         string
	Name        string
	Company     string
	Position    string
	Title       string
	Field       string
	Skill       string
	Mobile      string
	Email       string
	Projectname string
	Project     []map[string]string
	Content     string
	Templateid  string
}

//专家电子名片详情
func Expert(c echo.Context) error {

	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userid := m["userid"].(string)

	cn := pb.NewExpertClient(Connections["expert"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	r, err := cn.Info(ctx, &pb.InfoRequest{Userid: userid})
	if err != nil {
		Logger.Error("call micro service 'Expert' error: ", err)
		return errors.New("internal service error")
	}

	a := make(map[string]string)
	a["name"] = r.Name
	a["company"] = r.Company
	a["position"] = r.Position
	a["title"] = r.Title
	a["field"] = r.Field
	a["skill"] = r.Skill
	a["mobile"] = r.Mobile
	a["email"] = r.Email
	a["content"] = r.Content
	a["templateid"] = r.Templateid
	a["pic"] = r.Pic
	a["projectname"] = r.Projectname

	var b map[string]map[string]string
	json.Unmarshal([]byte(r.Project), &b)

	b["info"] = a

	return c.JSON(http.StatusOK, b)

}

//电子名片更新
func Fill(c echo.Context) error {
	// Use the following code piece for json body
	m := new(ExpertInfo)
	if err := c.Bind(&m); err != nil {
		return err
	}
	userid := m.Userid
	name := m.Name
	company := m.Company
	position := m.Position
	title := m.Title
	field := m.Field
	skill := m.Skill
	mobile := m.Mobile
	email := m.Email

	//校验token
	authtoken := AuthToken(c, mobile)
	if authtoken != "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result":  "false",
			"message": "token error",
		})
	}

	b, err := json.Marshal(m.Project)
	if err != nil {
		panic(err)
	}
	project := string(b)

	content := m.Content
	templateid := m.Templateid
	pic := m.Pic
	projectname := m.Projectname

	cn := pb.NewExpertClient(Connections["expert"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	r, err := cn.Fill(ctx, &pb.FillRequest{
		Userid:      userid,
		Name:        name,
		Company:     company,
		Position:    position,
		Title:       title,
		Field:       field,
		Skill:       skill,
		Mobile:      mobile,
		Email:       email,
		Project:     project,
		Content:     content,
		Templateid:  templateid,
		Pic:         pic,
		Projectname: projectname,
	})
	if err != nil {
		Logger.Error("call micro service 'Fill' error: ", err)
		return errors.New("internal service error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"filled": r.Filled,
	})

}

//电子名片收藏
func Collect(c echo.Context) error {
	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userid := m["userid"].(string)
	expertid := m["expertid"].(string)
	mobile := m["mobile"].(string)

	//校验token
	authtoken := AuthToken(c, mobile)
	if authtoken != "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result":  "false",
			"message": "token error",
		})
	}

	cn := pb.NewExpertClient(Connections["expert"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	_, err := cn.Collect(ctx, &pb.CollectRequest{Userid: userid, Expertid: expertid})
	if err != nil {
		Logger.Error("call micro service 'Collect' error: ", err)
		return errors.New("internal service error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"result": "true",
	})

}

//收藏列表
func CollectList(c echo.Context) error {
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

	cn := pb.NewExpertClient(Connections["expert"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	r, err := cn.CollectList(ctx, &pb.CollectListRequest{Userid: userid})
	if err != nil {
		Logger.Error("call micro service 'CollectList' error: ", err)
		return errors.New("internal service error")
	}

	return c.String(http.StatusOK, r.Collectlist)

}

//收藏删除
func CollectDel(c echo.Context) error {
	// Use the following code piece for json body
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	userid := m["userid"].(string)
	expertid := m["expertid"].(string)
	mobile := m["mobile"].(string)

	//校验token
	authtoken := AuthToken(c, mobile)
	if authtoken != "true" {
		return c.JSON(http.StatusOK, map[string]string{
			"result":  "false",
			"message": "token error",
		})
	}

	cn := pb.NewExpertClient(Connections["expert"])

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Attach rid to ctx, so that the micro service can access the same rid
	ctx = metadata.AppendToOutgoingContext(ctx, "rid", c.Response().Header().Get(echo.HeaderXRequestID))
	_, err := cn.CollectDel(ctx, &pb.CollectDelRequest{Userid: userid, Expertid: expertid})
	if err != nil {
		Logger.Error("call micro service 'CollectDel' error: ", err)
		return errors.New("internal service error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"result": "true",
	})

}
