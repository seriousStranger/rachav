package proxy

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/kopkapozla/rachav/database"
)

func GetProxyHandler(
	transport *http.Transport,
	upstreamAddr, fallbackAddr string,
	panelHandler http.Handler,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := TryPanelRequest(w, r, panelHandler)
		if ok {
			return
		}

		isAuth := checkAuth(r.Header.Get("Proxy-Authorization"))
		if r.Method != http.MethodConnect || !isAuth {
			target, _ := url.Parse("http://" + fallbackAddr)
			proxy := httputil.NewSingleHostReverseProxy(target)

			proxy.ServeHTTP(w, r)

			return
		}

		toNaive(w, r, transport, upstreamAddr)
	}
}

func checkAuth(authHeader string) bool {
	const prefix = "Basic "

	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(payload), ":", 2) //nolint:mnd

	users, err := database.Load()
	if err != nil {
		return false
	}

	pass, ok := users[pair[0]]
	if !ok {
		return false
	}

	if subtle.ConstantTimeCompare([]byte(pass), []byte(pair[1])) == 0 {
		return false
	}

	return true
}
