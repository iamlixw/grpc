package ctrl

import (
	"errors"
	"flag"
	"google.golang.org/grpc"

	wonaming "arctron.cn/arctron/arcplus/etcd"
)

var (
	services = map[string]*string{
		"user":   flag.String("user", "user", "service name"),
		"expert": flag.String("expert", "expert", "service name"),
		"qiniu":  flag.String("qiniu", "qiniu", "service name"),
		"member": flag.String("member", "member", "service name"),
		"other":  flag.String("other", "other", "service name"),
	}

	reg = flag.String("reg", "http://127.0.0.1:2379", "register address")

	// Connections contain the client connections to micro services
	Connections = make(map[string]*grpc.ClientConn)
)

// ConnSrv connects micro service used in gateway
func ConnSrv() error {
	flag.Parse()

	for name, addr := range services {
		r := wonaming.NewResolver(*addr)
		b := grpc.RoundRobin(r)

		conn, err := grpc.Dial(*reg, grpc.WithInsecure(), grpc.WithBalancer(b))
		if err != nil {
			Logger.Error(`connect to '%s:%s' service failed: %v`, name, addr, err)
			return errors.New("init services connection error")
		}
		Connections[name] = conn
	}
	return nil
}

// CloseConn closes all connects to services
func CloseConn() {
	for _, conn := range Connections {
		if conn != nil {
			conn.Close()
		}
	}
}
