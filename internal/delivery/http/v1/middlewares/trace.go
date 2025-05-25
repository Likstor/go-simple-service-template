package middleware

import (
	"net/http"
	"service/internal/pkg/common"

	"github.com/google/uuid"
)

func SetupTrace(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common.SetValueIntoContext(r.Context(), common.TRACE_KEY, uuid.New().String())

		handler.ServeHTTP(w, r)
	})
}