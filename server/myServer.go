package server

type MyServer struct {
	Ip                 string
	Ttl                int
	PullConsulInterval int
	ServiceName        string
	ConsulHost         string
	Port               string
	GinMode            GinHttpMode
}
