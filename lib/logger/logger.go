package logger

import (
	"fmt"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

func logWithColor(color, level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	fmt.Printf("%s[%s] [%s] %s%s\n",
		color,
		timestamp,
		level,
		message,
		colorReset,
	)
}

func Info(format string, args ...interface{}) {
	logWithColor(colorGreen, "INFO", format, args...)
}

func Success(format string, args ...interface{}) {
	logWithColor(colorGreen, "SUCCESS", format, args...)
}

func Warn(format string, args ...interface{}) {
	logWithColor(colorYellow, "WARN", format, args...)
}

func Error(format string, args ...interface{}) {
	logWithColor(colorRed, "ERROR", format, args...)
}
