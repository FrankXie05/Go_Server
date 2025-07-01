package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logDir = "logs"

// InitLogDirectory 保证日志目录存在
func InitLogDirectory() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}
}

// LogError 写入 error 日志
func LogError(traceID string, msg string) {
	InitLogDirectory()

	date := time.Now().Format("2006-01-02")
	logFile := fmt.Sprintf("%s/error-%s.log", logDir, date)
	writeToFile(logFile, "ERROR", traceID, msg)
}

// LogInfo 写入 info 日志（可选）
func LogInfo(traceID string, msg string) {
	InitLogDirectory()

	date := time.Now().Format("2006-01-02")
	logFile := fmt.Sprintf("%s/info-%s.log", logDir, date)
	writeToFile(logFile, "INFO", traceID, msg)
}

func writeToFile(filename, level, traceID, msg string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("日志写入失败:", err)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("[%s] [%s] %s", level, traceID, msg)
}
