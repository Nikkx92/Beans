package logger

import (
	"log/slog"
	"path/filepath"
	"voicer/internal/storage"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Message struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"msg"`
	Source  struct {
		Function string `json:"function"`
		File     string `json:"file"`
		Line     int    `json:"line"`
	} `json:"source,omitempty"`
}

func NewLogger() {
	dir, err := storage.Path()
	if err != nil {
		slog.Info(err.Error())
	}

	writer := &lumberjack.Logger{
		Filename:   filepath.Join(dir, "app.log"),
		MaxSize:    10, // мегабайт
		MaxBackups: 3,
		MaxAge:     365,  // дней хранения
		Compress:   true, // сжатие старых логов gzip
	}
	defer writer.Close()

	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)
}
