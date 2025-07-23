package config

type (
	SSL struct {
		Timeout int `validate:"gte=0"` // In Millisecond
	}
)

var Default = &SSL{
	Timeout: 2000,
}
