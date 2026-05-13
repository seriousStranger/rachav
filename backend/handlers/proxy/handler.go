package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kopkapozla/rachav/auth"
	"github.com/kopkapozla/rachav/database"
)

func GetProxyHandler(
	transport *http.Transport,
	upstreamAddr, fallbackAddr string,
	panelHandler http.Handler,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header)
		ok := TryPanelRequest(w, r, panelHandler)
		if ok {
			return
		}

		users, err := database.Load()
		isAuth := false
		if err == nil {
			isAuth = auth.СheckAuth(r.Header.Get("Proxy-Authorization"), "Basic ", users)
		}
		if r.Method != http.MethodConnect || !isAuth {
			target, _ := url.Parse("http://" + fallbackAddr)
			proxy := httputil.NewSingleHostReverseProxy(target)

			proxy.ServeHTTP(w, r)

			return
		}

		toNaive(w, r, transport, upstreamAddr)
	}
}
