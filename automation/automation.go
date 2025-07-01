package automation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RunBrowserTask(scriptUUID string, tracklink string) {
	api := "http://localhost:36912/browser/open-url"
	data := map[string]interface{}{
		"scriptUUID": scriptUUID,
		"url":        tracklink,
	}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", api, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to open url:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(">>> OpenURL:", tracklink)
	fmt.Println(">>> scriptUUID:", scriptUUID)
	fmt.Println("OpenURL status:", resp.Status)
}
