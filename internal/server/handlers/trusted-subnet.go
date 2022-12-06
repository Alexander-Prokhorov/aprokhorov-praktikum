package handlers

import (
	"net"
	"net/http"
)

const xRealIP = "X-Real-IP"

func CheckACL(acl *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if acl == nil {
				next.ServeHTTP(w, r)

				return
			}

			realIPHeader := r.Header.Get(xRealIP)
			realIP := net.ParseIP(realIPHeader)
			// no check for nil, because it doesn't matter

			if acl.Contains(realIP) {
				next.ServeHTTP(w, r)
			}

			http.Error(w, "Only Trusted Net Allowed", http.StatusForbidden)
		})
	}
}
