package gologin

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// 执行 HTTP 请求并校验状态码、返回 body（适用于 JSON 响应或需回显内容）
func DoRequestWithCheck(req *http.Request, expectedStatusCodes ...int) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// 若未指定 expected 状态码，默认检查 isSuccessStatus
	if len(expectedStatusCodes) == 0 {
		if !isSuccessStatus(resp.StatusCode) {
			return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
	} else {
		found := false
		for _, code := range expectedStatusCodes {
			if resp.StatusCode == code {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
	}

	return bodyBytes, nil
}
