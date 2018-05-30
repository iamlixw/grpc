package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	wonaming "arctron.cn/arctron/arcplus/etcd"
	. "arctron.cn/arctron/arcplus/other/internal"
	pb "arctron.cn/arctron/arcplus/pb/other"
)

const (
	name    = "other" // service name
	defPort = 13005   // default listening port
)

var (
	serv = flag.String("other", "other", "service name")
	port = flag.Int("port", defPort, "listening port")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register address")
)

func main() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = wonaming.Register(*serv, "127.0.0.1", *port, *reg, time.Second*10, 15)
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		log.Fatalf("receive signal '%v'", s)
		wonaming.UnRegister()
		os.Exit(1)
	}()

	Logger.Infof("service %s is starting at port %d", name, *port)

	s := grpc.NewServer()
	pb.RegisterOtherServer(s, &OtherServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
