package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type ShortSourceHandler struct {
	slog.Handler
}

func (h ShortSourceHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			r.AddAttrs(slog.String("source", filepath.Base(f.File)+":"+strconv.Itoa(f.Line)))
		}
	}
	return h.Handler.Handle(ctx, r)
}
func SetupSlog(level slog.Leveler) {
	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(ShortSourceHandler{textHandler}))
}
