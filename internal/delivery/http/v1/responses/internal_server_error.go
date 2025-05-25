package responses

import (
	"context"
	"net/http"
)

func InternalServerError(ctx context.Context, w http.ResponseWriter) {
	Error(
		ctx,
		w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
	)
}