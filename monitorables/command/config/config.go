package config

type Command struct {
	Timeout int `validate:"gte=0"` // In Millisecond
}

var Default = &Command{
	Timeout: 10000,
}
