package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/biisal/db-gui/internal/color"
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
func SetupSlog(filePath string, level slog.Leveler) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			color.Default.Info("The filepath `%s` doesn't exits\nDo you want to create it?[y/n]", filePath)
			var response string
			_, scanErr := fmt.Scanln(&response)
			if scanErr != nil || (response != "y" && response != "Y") {
				color.Default.Error("Log file creation aborted. Exiting...")
				return nil, err
			}
			dir := filepath.Dir(filePath)
			if mkErr := os.MkdirAll(dir, os.ModePerm); mkErr != nil {
				slog.Error("Failed to create log directory", "error", mkErr)
				return nil, mkErr
			}
			file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				slog.Error("Failed to create log file", "error", err)
				return nil, err
			}

		} else {
			slog.Error("Failed to open log file", "error", err)
			return nil, err
		}
	}
	textHandler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(ShortSourceHandler{textHandler}))
	return file, nil
}
