package internal

import (
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"

	. "arctron.cn/arctron/arcplus/lib"
	pb "arctron.cn/arctron/arcplus/pb/qiniu"
)

var (
	Logger  = log.New("qiniu")
	Bucket  = "demo"
	AK      = "demoak"
	SK      = "demosk"
	upToken string
)

type QiniuServer struct{}

//GetToken 获取token
func (u *QiniuServer) GetToken(c context.Context, aq *pb.UpTokenRequest) (*pb.UpTokenReply, error) {
	Logger.Infof("request id = %s", GetRid(c))

	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AK, SK)
	upToken = putPolicy.UploadToken(mac)

	return &pb.UpTokenReply{Token: upToken}, nil
}
