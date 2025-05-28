package responses

import (
	"context"
	"net/http"
)

func Error(ctx context.Context, w http.ResponseWriter, statusCode int, msg string) {
	resp := make(map[string]any)
	resp["error"] = msg

	JSON(ctx, w, statusCode, resp)
}