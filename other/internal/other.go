package internal

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"

	. "arctron.cn/arctron/arcplus/lib"
	pb "arctron.cn/arctron/arcplus/pb/other"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Logger = log.New("other")
	db     = &sql.DB{}
	dbuser string
	dbpass string
)

type OtherServer struct{}

//发送短信验证码
func (u *OtherServer) SendSms(c context.Context, aq *pb.SendRequest) (*pb.SendReply, error) {
	Logger.Infof("request id = %s", GetRid(c))

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	var s = NewAliSms()
	err := s.Send(aq.Mobile, `{"code":"`+vcode+`"}`, "SMS_122289786")

	if err != nil {
		return &pb.SendReply{Sended: "false"}, nil
	}

	t := time.Now()

	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
	db, _ = sql.Open("mysql", dbuser+":"+dbpass+"@/arcplus")
	db.Exec(
		"INSERT INTO arc_sms (mobile,code,createTime) VALUES (?,?,?)",
		aq.Mobile,
		vcode,
		t.Unix(),
	)

	defer db.Close()

	return &pb.SendReply{Sended: "true"}, nil
}

//校验验证码
func (u *OtherServer) CheckCode(c context.Context, aq *pb.CheckRequest) (*pb.CheckReply, error) {
	Logger.Infof("request id = %s", GetRid(c))

	t := time.Now().Unix() - 300

	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
	db, _ = sql.Open("mysql", dbuser+":"+dbpass+"@/arcplus")
	rows := db.QueryRow("select code from arc_sms where mobile = ? and createTime > ? order by createTime desc limit 1", aq.Mobile, t)
	var code string
	rows.Scan(&code)

	defer db.Close()

	if code == aq.Code {
		return &pb.CheckReply{Checked: "true"}, nil
	} else {
		return &pb.CheckReply{Checked: "false"}, nil
	}

}
