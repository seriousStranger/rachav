package proxy_handler

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kopkapozla/rachav/config"
	"github.com/kopkapozla/rachav/database"
	"github.com/labstack/echo/v5"
)

type byteCounter struct {
	writer io.Writer
	count  atomic.Int64
	label  string
}

func (bc *byteCounter) Write(p []byte) (int, error) {
	n, err := bc.writer.Write(p)
	if n > 0 {
		total := bc.count.Add(int64(n))
		if total%102400 < int64(n) {
			log.Printf("[%s] Передано: %.2f КБ", bc.label, float64(total)/1024)
		}
	}
	return n, err
}

func GetProxyHandler(tr *http.Transport, upstreamAddr, fallbackAddr string, echoHandlers *echo.Echo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.URL.Path)

		if config.Config.IsPanelEnable() &&
			strings.HasPrefix(r.URL.Path, "/"+config.Config.GetPanelUrl()) {
			echoHandlers.ServeHTTP(w, r)
		}

		isAuth, _ := checkAuth(r.Header.Get("Proxy-Authorization"))
		if r.Method != http.MethodConnect || !isAuth {
			target, _ := url.Parse("http://" + fallbackAddr)

			// Initialize the proxy
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, r)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Hour)
		defer cancel()

		rc := http.NewResponseController(w)

		pr, pw := io.Pipe()
		outReq, _ := http.NewRequestWithContext(ctx, http.MethodConnect, "http://"+upstreamAddr, pr)
		outReq.Host = r.Host

		resp, err := tr.RoundTrip(outReq)
		if err != nil {
			log.Printf("Upstream error (%s): %v", r.Host, err)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		rc.Flush()

		toUpstream := &byteCounter{writer: pw, label: "C -> S"}
		toClient := &byteCounter{writer: w, label: "S -> C"}

		errChan := make(chan error, 2)

		go func() {
			buf := make([]byte, 32*1024)
			for {
				n, rerr := r.Body.Read(buf)
				if n > 0 {
					_, werr := toUpstream.Write(buf[:n])
					if werr != nil {
						errChan <- werr
						return
					}
				}
				if rerr != nil {
					pw.Close()
					if rerr != io.EOF {
						errChan <- rerr
					} else {
						errChan <- nil
					}
					return
				}
			}
		}()

		go func() {
			buf := make([]byte, 32*1024)
			for {
				n, rerr := resp.Body.Read(buf)
				if n > 0 {
					_, werr := toClient.Write(buf[:n])
					if werr != nil {
						errChan <- werr
						return
					}
					_ = rc.Flush()
				}
				if rerr != nil {
					if rerr != io.EOF {
						errChan <- rerr
					} else {
						errChan <- nil
					}
					return
				}
			}
		}()

		err = <-errChan
		if err != nil && err != context.Canceled {
			log.Printf("Stream finished with error: %v", err)
		}
		log.Printf("Session %s ended", r.Host)
	}
}

func checkAuth(authHeader string) (bool, error) {
	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return false, fmt.Errorf("invalid prefix")
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return false, err
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return false, fmt.Errorf("invalid format")
	}

	users, err := database.Load()
	if err != nil {
		return false, err
	}

	pass, ok := users[pair[0]]
	if !ok {
		return false, nil
	}
	if subtle.ConstantTimeCompare([]byte(pass), []byte(pair[1])) == 0 {
		return false, nil
	}

	return true, nil
}
