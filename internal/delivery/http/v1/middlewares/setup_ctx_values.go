package middleware

import (
	"service/internal/pkg/common"
	"context"
	"net/http"
)

func SetupContextValues(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valuesMap := make(map[string]any)
		
		r = r.WithContext(context.WithValue(r.Context(), common.CTX_VALUES_KEY, valuesMap))

		handler.ServeHTTP(w, r)
	})
}