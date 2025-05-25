package logs

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"service/internal/pkg/common"
	"time"
)

var ErrLoggerSetup = errors.New("logger setup error")

var DEBUG = false

func init() {
	DEBUG = os.Getenv("ENV") == "DEBUG"
}

type PrettyWriter struct {
	writer io.Writer
}

func (pw PrettyWriter) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, p, "", "  "); err != nil {
		return pw.writer.Write(p)
	}
	return pw.writer.Write(buf.Bytes())
}

type CustomHandler struct {
	slog.Handler
}

func (h CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if common.CTXValuesExists(ctx) {
		trace := common.GetTrace(ctx)

		if trace != "" {
			r.AddAttrs(slog.String("trace", trace))
		}
	}

	if DEBUG {
		var pcs [1]uintptr
		runtime.Callers(5, pcs[:])

		r.AddAttrs(
			getSource(pcs[0]),
		)
	}

	return h.Handler.Handle(ctx, r)
}

func getSource(pc uintptr) slog.Attr {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()

	return slog.Group(
		"source",
		slog.String("func", f.Function),
		slog.Int("line", f.Line),
	)
}

var op = "logs.Setup"

func Setup(ctx context.Context, logsDir string) (*os.File, *bufio.Writer, error) {
	abs, err := filepath.Abs(logsDir)
	if err != nil {
		Error(
			ctx,
			err.Error(),
			op,
		)
	}

	if err := os.MkdirAll(abs, os.ModeDir); err != nil {
		Error(
			ctx,
			err.Error(),
			op,
		)

		return nil, nil, ErrLoggerSetup
	}

	file, err := os.OpenFile(filepath.Join(abs, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02_15-04-05Z07-00"))), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		Error(
			ctx,
			err.Error(),
			op,
		)

		return nil, nil, ErrLoggerSetup
	}

	writer := bufio.NewWriter(io.MultiWriter(os.Stdout, file))

	var writerForLogs io.Writer

	if os.Getenv("ENV") == "DEBUG" {
		writerForLogs = PrettyWriter{writer: writer}
	} else {
		writerForLogs = writer
	}

	handler := CustomHandler{
		Handler: slog.NewJSONHandler(writerForLogs, nil),
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	go func() {
		for {
			time.Sleep(time.Minute * 1)
			writer.Flush()
		}
	}()

	return file, writer, nil
}

func Error(ctx context.Context, msg, op string, other ...any) {
	slog.ErrorContext(
		ctx,
		msg,
		slog.String("op", op),
		slog.Group("data", other...),
	)
}

func Info(ctx context.Context, msg, op string, other ...any) {
	slog.InfoContext(
		ctx,
		msg,
		slog.String("op", op),
		slog.Group("data", other...),
	)
}

func Warn(ctx context.Context, msg, op string, other ...any) {
	slog.WarnContext(
		ctx,
		msg,
		slog.String("op", op),
		slog.Group("data", other...),
	)
}
