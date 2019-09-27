package consulStruct

import "time"

type ServerInfo struct {
	ServiceName string
	Ip          string
	Port        string
	Desc        string
	UpdateTime  time.Time
	CreateTime  time.Time
	Ttl         int
	ServerType  int // 1 http 2 grpc
}

const (
	HttpType = 1
	GrpcType = 2
)
