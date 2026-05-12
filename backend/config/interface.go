package config

type ConfigInterface interface {
	Init()
	GetAuthPair() (string, string)

	IsPanelEnable() bool
	GetPanelUrl() string

	GetListenPort() string
	GetFallbackPort() string
	GetNaivePort() string

	GetHost() string
}

var Config ConfigInterface //nolint

func SetConfigImplementation(implementation ConfigInterface) {
	Config = implementation
}
