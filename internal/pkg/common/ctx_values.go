package common

import (
	"context"
	"service/internal/domain/types"
)

const (
	CTX_VALUES_KEY = types.ContextValueKey("ctx_values")
	TRACE_KEY = "trace"
)

func CTXValuesExists(ctx context.Context) bool {
	switch ctx.Value(CTX_VALUES_KEY).(type) {
	case map[string]any:
		return true
	default:
		return false
	}
}

func GetValueFromContext(ctx context.Context, key string) any {
	values := ctx.Value(CTX_VALUES_KEY).(map[string]any)

	return values[key]
}

func SetValueIntoContext(ctx context.Context, key string, value any) {
	values := ctx.Value(CTX_VALUES_KEY).(map[string]any)

	values[key] = value
}

func GetTrace(ctx context.Context) string {
	correlationID := GetValueFromContext(ctx, TRACE_KEY)
	
	switch v := correlationID.(type) {
	case string:
		return v
	}

	return ""
}