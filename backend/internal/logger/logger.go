package logger

import (
	"log"
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *slog.Logger

func InitLogger() {
	fileHandler := &lumberjack.Logger{
		Filename:   "logs/app.log", // Путь к файлу логов
		MaxSize:    10,             // Макс размер файла (MB)
		MaxBackups: 2,              // Количество резервных копий
		MaxAge:     1,              // Дни хранения
		Compress:   true,           // Сжатие файлов логов
	}

	jsonHandler := slog.NewJSONHandler(fileHandler, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Логирование с уровня Info
	})

	Log = slog.New(jsonHandler)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Fatal(message string, err error) {
	Log.Error(message, slog.String("error", err.Error()))
	log.Fatal(err)
}
