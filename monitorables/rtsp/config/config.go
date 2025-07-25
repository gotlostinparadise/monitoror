package config

type (
	RTSP struct {
		Timeout int `validate:"gte=0"`
	}
)

var Default = &RTSP{
	Timeout: 2000,
}
