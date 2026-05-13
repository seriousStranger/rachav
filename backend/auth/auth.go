package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

func СheckAuth(authHeader, prefix string, users map[string]string) bool {
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(payload), ":", 2) //nolint:mnd

	pass, ok := users[pair[0]]

	if subtle.ConstantTimeCompare([]byte(pass), []byte(pair[1])) == 0 {
		return false
	}

	if !ok {
		return false
	}

	return true
}
