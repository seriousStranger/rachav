package proxy

import (
	"net/http"
	"strings"

	"github.com/kopkapozla/rachav/config"
)

func TryPanelRequest(w http.ResponseWriter, r *http.Request, panelHandler http.Handler) bool {
	if config.Config.IsPanelEnable() {
		path, found := strings.CutPrefix(r.URL.Path, "/"+config.Config.GetPanelUrl())
		if found {
			r.URL.Path = path
			panelHandler.ServeHTTP(w, r)

			return true
		}
	}
	return false
}
