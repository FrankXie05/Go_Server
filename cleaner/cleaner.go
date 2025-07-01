package cleaner

import (
	"EMU_server/gologin"
	"EMU_server/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCleanupScheduler(gl *gologin.GoLogin, maxProfiles int, whitelist map[string]bool, cronSpec string, webhookURL string) {
	fmt.Println("Cleanup scheduler started with cron spec:", cronSpec)
	c := cron.New()
	c.AddFunc(cronSpec, func() {
		traceID := fmt.Sprintf("cleanup-%d", time.Now().UnixNano())
		summary, err := gologin.CleanupExcessProfiles(gl, maxProfiles, whitelist, traceID)

		if err != nil {
			logger.LogError(traceID, err.Error())
		} else {
			logger.LogInfo(traceID, summary)
			sendWebhook(traceID, webhookURL, summary)
		}
	})
	c.Start()
}

func sendWebhook(traceID, url, summary string) {
	payload := map[string]interface{}{
		"traceId": traceID,
		"message": summary,
		"source":  "auto-profile-cleaner",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.LogError(traceID, "Webhook POST 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()
	logger.LogInfo(traceID, fmt.Sprintf("Webhook 已发送，状态码: %d", resp.StatusCode))
}
