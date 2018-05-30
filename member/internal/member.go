package internal

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"

	. "arctron.cn/arctron/arcplus/lib"
	pb "arctron.cn/arctron/arcplus/pb/member"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Logger = log.New("member")
	db     = &sql.DB{}
	dbuser string
	dbpass string
)

type MemberServer struct{}

//激活
func (u *MemberServer) Regist(c context.Context, aq *pb.RegistRequest) (*pb.RegistReply, error) {
	Logger.Infof("request id = %s", GetRid(c))

	id := uuid.New()
	t := time.Now()

	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
	db, _ = sql.Open("mysql", dbuser+":"+dbpass+"@/arcplus")

	rows := db.QueryRow("select id as userid from arc_user where mobile = ?", aq.Mobile)
	var userid string
	rows.Scan(&userid)

	if userid == "" {
		db.Exec(
			"INSERT INTO arc_user (id,username,name,mobile,email,createTime) VALUES (?,?,?,?,?,?)",
			id.String(),
			aq.Mobile,
			aq.Name,
			aq.Mobile,
			aq.Email,
			t.Unix(),
		)
	}

	defer db.Close()

	return &pb.RegistReply{Registed: "true"}, nil
}
