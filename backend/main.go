package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/kopkapozla/rachav/auth"
	"github.com/kopkapozla/rachav/config"
	"github.com/kopkapozla/rachav/config/viper"
	"github.com/kopkapozla/rachav/database"
	"github.com/kopkapozla/rachav/handlers/panel"
	proxy_handler "github.com/kopkapozla/rachav/handlers/proxy"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

const (
	IdleTimeout          = 5 * time.Minute
	MaxConcurrentStreams = 1000
	TLSHandshakeTimeout  = 5 * time.Minute
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
		_, err := strconv.Atoi(config.Config.GetNaivePort())
		if err != nil {
			panic("strange port value, try to trim it")
		}

		// https://github.com/seriousStranger/rachav/issues/1
		//nolint:gosec
		cmd := exec.CommandContext(
			context.Background(),
			"./naive",
			"--listen=http://127.0.0.1:"+config.Config.GetNaivePort(),
		)

		err = cmd.Run()
		if err != nil {
			slog.Error(err.Error())
			panic(err)
		}
	}()

	certmanager := &autocert.Manager{
		Cache:      autocert.DirCache("./certs"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Config.GetHost()),
	}

	transport := &http.Transport{
		IdleConnTimeout:     IdleTimeout,
		TLSHandshakeTimeout: TLSHandshakeTimeout,
		WriteBufferSize:     0,
		ReadBufferSize:      0,
	}
	transport.Protocols = new(http.Protocols)
	transport.Protocols.SetHTTP1(true)
	transport.Protocols.SetUnencryptedHTTP2(true)

	var echoServer *echo.Echo

	if config.Config.IsPanelEnable() {
		slog.Warn("panel enable")

		echoServer = echo.New()

		echoServer.GET("/", getPanelHtml(panelFile))

		api := echoServer.Group("/api")
		api.Use(auth.EchoMiddleware)

		api.GET("/user-list", panel.GetUserList)
		api.POST("/user-list", panel.PostUserList)
		api.GET("/host", panel.GetHost)
	}

	handler := http.HandlerFunc(
		proxy_handler.GetProxyHandler(
			transport,
			"127.0.0.1:"+config.Config.GetNaivePort(),
			"127.0.0.1:"+config.Config.GetFallbackPort(),
			echoServer,
		),
	)

	server := &http.Server{
		Addr:         ":" + config.Config.GetListenPort(),
		Handler:      handler,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  IdleTimeout,
		TLSConfig:    certmanager.TLSConfig(),
	}

	server.Protocols = new(http.Protocols)
	server.Protocols.SetHTTP1(true)
	server.Protocols.SetHTTP2(true)

	http2server := &http2.Server{
		MaxConcurrentStreams: MaxConcurrentStreams,
		IdleTimeout:          IdleTimeout,
	}

	err = http2.ConfigureServer(server, http2server)
	if err != nil {
		panic(err)
	}

	slog.Info("racahv on " + config.Config.GetHost() + ":" + config.Config.GetListenPort())
	slog.Error(server.ListenAndServeTLS("", "").Error())
}

func getPanelHtml(files fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		html, err := fs.ReadFile(files, "build_frontend/index.html")
		if err != nil {
			slog.Error(
				"Can't read index.html. In base configuration it's embed, so in 99% time it's code problem",
			)

			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"something is broken: "+err.Error(),
			)
		}

		return c.HTMLBlob(http.StatusOK, html)
	}
}
