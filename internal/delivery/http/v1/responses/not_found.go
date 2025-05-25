package responses

import (
	"context"
	"net/http"
)

func NotFound(ctx context.Context, w http.ResponseWriter) {
	Error(
		ctx,
		w,
		http.StatusNotFound,
		http.StatusText(http.StatusNotFound),
	)
}