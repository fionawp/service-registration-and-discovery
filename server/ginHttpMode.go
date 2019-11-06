package server

//gin http_mode
type GinHttpMode int

const (
	DebugMode GinHttpMode = iota
	ReleaseMode
	TestMode
)

func (hm GinHttpMode) String() string {
	switch hm {
	case DebugMode:
		return "debug"
	case ReleaseMode:
		return "release"
	case TestMode:
		return "test"
	default:
		return ""
	}
}
