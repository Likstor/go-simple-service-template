package middleware

import (
	"net/http"
	"service/internal/pkg/common"

	"github.com/google/uuid"
)

func SetupTrace(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := common.GetRequestID(r); reqID == "" {
			r.Header.Set("X-Request-ID", uuid.NewString())
		}

		common.SetValueIntoContext(r.Context(), common.TRACE_KEY, r.Header.Get("X-Request-ID"))

		handler.ServeHTTP(w, r)
	})
}