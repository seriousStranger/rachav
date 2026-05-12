package main

import (
	"crypto/subtle"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os/exec"
	"time"

	"github.com/kopkapozla/rachav/config"
	"github.com/kopkapozla/rachav/config/viper"
	"github.com/kopkapozla/rachav/database"
	"github.com/kopkapozla/rachav/handlers"
	"github.com/kopkapozla/rachav/proxy_handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

//go:embed build_frontend
var panelFile embed.FS

func main() {
	err := database.CreateDbIfNotExist()
	if err != nil {
		panic(err)
	}

	config.SetConfigImplementation(viper.NewViperConfig())
	config.Config.Init()

	go func() {
		cmd := exec.Command("./naive", "--listen=http://127.0.0.1:"+config.Config.GetNaivePort())
		err := cmd.Run()
		if err != nil {
			fmt.Println("error:", err)
			panic(err)
		}
	}()

	m := &autocert.Manager{
		Cache:      autocert.DirCache("./certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Config.GetHost()),
	}

	tr := &http.Transport{
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		WriteBufferSize:     0,
		ReadBufferSize:      0,
	}
	tr.Protocols = new(http.Protocols)
	tr.Protocols.SetHTTP1(true)
	tr.Protocols.SetUnencryptedHTTP2(true)

	echoServer := echo.New()
	if config.Config.IsPanelEnable() {
		slog.Warn("panel enable")

		panel := echoServer.Group("/" + config.Config.GetPanelUrl())

		panel.GET("/", getPanelHtml(panelFile))

		api := panel.Group("/api")
		api.Use(middleware.BasicAuth(authForApi))

		api.GET("/user-list", handlers.GetUserList)
		api.POST("/user-list", handlers.PostUserList)
	}

	handler := http.HandlerFunc(
		proxy_handler.GetProxyHandler(
			tr,
			"127.0.0.1:"+config.Config.GetNaivePort(),
			"127.0.0.1:"+config.Config.GetFallbackPort(),
			echoServer,
		),
	)

	srv := &http.Server{
		Addr:         ":" + config.Config.GetListenPort(),
		Handler:      handler,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  10 * time.Minute,
		TLSConfig:    m.TLSConfig(),
	}

	srv.Protocols = new(http.Protocols)
	srv.Protocols.SetHTTP1(true)
	srv.Protocols.SetHTTP2(true)

	h2s := &http2.Server{
		MaxConcurrentStreams: 1000,
		IdleTimeout:          5 * time.Minute,
	}

	http2.ConfigureServer(srv, h2s)

	slog.Info("racahv on " + config.Config.GetHost() + ":" + config.Config.GetListenPort())
	slog.Error(srv.ListenAndServeTLS("", "").Error())
}

func authForApi(c *echo.Context, user string, password string) (bool, error) {
	curUser, curPass := config.Config.GetAuthPair()
	if subtle.ConstantTimeCompare([]byte(user), []byte(curUser)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(curPass)) == 1 {
		return true, nil
	}
	slog.Warn(user + ":" + password)
	slog.Warn(curUser + ":" + curPass)
	return false, nil
}

func getPanelHtml(files fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		html, err := fs.ReadFile(files, "build_frontend/index.html")
		if err != nil {
			return echo.NewHTTPError(500, "something went wrong: "+err.Error())
		}

		return c.HTMLBlob(200, html)
	}
}
