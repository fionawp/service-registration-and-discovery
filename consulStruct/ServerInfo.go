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
}
