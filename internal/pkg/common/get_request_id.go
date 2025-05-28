package common

import "net/http"

func GetRequestID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}