package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

var Gl *slog.Logger

func Initialize() {

	file, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	w := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(w, opts)

	Gl = slog.New(handler)

	slog.SetDefault(Gl)
}
