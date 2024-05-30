package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init(loggerLevel string) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(config.EncoderConfig)

	logFile, err := os.OpenFile("C:/CoggersProject/logs/golog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка создания файла golog %s", err)
	}
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	var fileLogLevel zapcore.Level
	var consoleLogLevel zapcore.Level

	switch loggerLevel {
	case "debug":
		fileLogLevel = zapcore.InfoLevel
		consoleLogLevel = zapcore.InfoLevel
	default:
		fileLogLevel = zapcore.DPanicLevel
		consoleLogLevel = zapcore.DebugLevel
	}

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, fileLogLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, consoleLogLevel)

	core := zapcore.NewTee(fileCore, consoleCore)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
