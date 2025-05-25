package responses

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"service/internal/pkg/logs"
)

const opJSON = "responses.JSON"

// JSON формирует ответ в формате JSON и записывает его в http.ResponseWriter.
func JSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logs.Error(
			ctx,
			err.Error(),
			opJSON,
		)

		InternalServerError(ctx, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		logs.Error(
			ctx,
			err.Error(),
			opJSON,
		)

		return
	}

	logs.Info(
		ctx,
		"written answer",
		opJSON,
		slog.Int("statusCode", statusCode),
	)
}