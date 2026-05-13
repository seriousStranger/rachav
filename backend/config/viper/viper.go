package viper

import (
	"crypto/rand"
	"log/slog"

	"github.com/spf13/viper"
)

const (
	defaultLogin          = "user"
	defaultPassword       = "ChangeThisPasswordPlease"
	enablePanel           = false
	defaultListenPort     = 443
	defaultNaiveProxyPort = 8080
	defaultFallbackPort   = 8081
	defaultHost           = "localhost"
)

type ViperConfig struct{}

func (vp ViperConfig) Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err == nil {
		slog.Any("config file", slog.String("path", viper.ConfigFileUsed()))

		return
	}

	slog.Warn("can't read config...")

	defaultPanelUrl := rand.Text()

	viper.SetDefault("login", defaultLogin)
	viper.SetDefault("password", defaultPassword)
	viper.SetDefault("enable_panel", enablePanel)
	viper.SetDefault("panel_url", defaultPanelUrl)
	viper.SetDefault("listen_port", defaultListenPort)
	viper.SetDefault("naiveproxy_port", defaultNaiveProxyPort)
	viper.SetDefault("fallback_port", defaultFallbackPort)
	viper.SetDefault("host", defaultHost)

	slog.Warn("creating example config... Change it as fast as you can :)")

	err = viper.SafeWriteConfigAs("config.yaml")
	if err != nil {
		slog.Warn("could not create default config", "error", err)
	}
}

func (vp ViperConfig) GetAuthPair() (string, string) {
	user := viper.GetString("login")
	password := viper.GetString("password")

	return user, password
}

func (vp ViperConfig) IsPanelEnable() bool {
	return viper.GetBool("enable_panel")
}

func (vp ViperConfig) GetPanelUrl() string {
	return viper.GetString("panel_url")
}

func (vp ViperConfig) GetListenPort() string {
	return viper.GetString("listen_port")
}

func (vp ViperConfig) GetFallbackPort() string {
	return viper.GetString("fallback_port")
}

func (vp ViperConfig) GetNaivePort() string {
	return viper.GetString("naiveproxy_port")
}

func (vp ViperConfig) GetHost() string {
	return viper.GetString("host")
}

func NewViperConfig() *ViperConfig {
	return &ViperConfig{}
}
