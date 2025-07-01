package main

import (
	"EMU_server/cleaner"
	"EMU_server/gologin"
	"EMU_server/server"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// 启动主服务
	r := server.SetupRouter()
	r.Run(":8090")
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gl := gologin.NewGoLogin()

	whitelist := map[string]bool{
		"": true,
	}

	// 启动定时清理器（每 10 分钟运行，cron 表达式）
	cleanupCron := "*/10 * * * *" // 每 10 分钟执行一次
	webhookURL := os.Getenv("CLEANUP_WEBHOOK_URL")
	go cleaner.StartCleanupScheduler(gl, 100, whitelist, cleanupCron, webhookURL)

	localPaths := []string{
		`C:\Users\Administrator\AppData\Local\Temp\2`,
		`C:\Users\Administrator\AppData\Local\Temp\GoLogin\profiles`,
	}
	cleaner.StartLocalCleanupForever(2*time.Hour, localPaths)
}
