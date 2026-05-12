package viper

import (
	"crypto/rand"
	"log/slog"

	"github.com/spf13/viper"
)

type ViperConfig struct{}

func (vp ViperConfig) Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err == nil {
		return
	}

	viper.SetDefault("login", "user")
	viper.SetDefault("password", "ChangeThisPasswordPlease")
	viper.SetDefault("Enable_panel", false)
	viper.SetDefault("panel_url", rand.Text())
	viper.SetDefault("listen_port", 443)
	viper.SetDefault("naiveproxy_port", 8080)
	viper.SetDefault("fallback_port", 8081)
	viper.SetDefault("host", "localhost")

	slog.Info("creating example config... Change it as fast as you can)")

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
