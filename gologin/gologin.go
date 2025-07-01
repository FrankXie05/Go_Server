package gologin

import (
	"EMU_server/logger"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GoLogin 结构体
type GoLogin struct {
	Token string
}

// 初始化 GoLogin
func NewGoLogin() *GoLogin {
	token := os.Getenv("GO_LOGIN_TOKEN")
	if token == "" {
		panic(" GO_LOGIN_TOKEN 未设置")
	}
	return &GoLogin{Token: token}
}

// 创建 GoLogin 指纹浏览器环境
func CreateProfile(gl *GoLogin, config *ProfileConfig) (string, error) {
	url := "https://api.gologin.com/browser"
	body, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshal config failed: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("build request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+gl.Token)
	req.Header.Set("Content-Type", "application/json")

	respBody, err := DoRequestWithCheck(req, http.StatusCreated, http.StatusOK)
	if err != nil {
		return "", fmt.Errorf("create profile failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode response failed: %v", err)
	}

	id, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response: missing id")
	}
	return id, nil
}

// 启动 GoLogin Profile
func StartProfile(gl *GoLogin, profileID string) (string, error) {
	url := "http://localhost:36912/browser/start-profile"
	payload := map[string]interface{}{
		"profileId": profileID,
		"sync":      true,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+gl.Token)
	req.Header.Set("Content-Type", "application/json")

	respBody, err := DoRequestWithCheck(req)
	if err != nil {
		return "", fmt.Errorf("start profile failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode response: %v", err)
	}

	addr, ok := result["wsUrl"].(string)
	if !ok {
		return "", fmt.Errorf("no wsUrl in response: %v", result)
	}
	return addr, nil
}

func CleanCache(profileID string) {
	fmt.Println("******** Cleaning Profile Cache ********", profileID)
	profilePath := filepath.Join(os.TempDir(), "Gologin", "profiles", profileID)
	if err := os.RemoveAll(profilePath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("无法删除目录 %s：%v\n", profilePath, err)
	} else {
		fmt.Printf("已清理: %s\n", profilePath)
	}

	// artifactBase := "C:\\Users\\Administrator\\AppData\\Local\\Temp"
	// art_subdir := []string{"2"}
	// for _, sub := range art_subdir {
	// 	path := filepath.Join(artifactBase, sub)
	// 	entries, err := os.ReadDir(path)
	// 	if err != nil {
	// 		fmt.Printf("读取目录失败: %s (%v)\n", path, err)
	// 		continue
	// 	}
	// 	for _, entry := range entries {
	// 		if strings.HasPrefix(entry.Name(), "playwright-artifacts-") {
	// 			target := filepath.Join(path, entry.Name())
	// 			err := os.RemoveAll(target)
	// 			if err != nil {
	// 				fmt.Printf("删除失败: %s (%v)\n", target, err)
	// 			} else {
	// 				fmt.Printf("已清理: %s\n", target)
	// 			}
	// 		}
	// 	}
	// }

}

func StopProfile(debuggerAddress string) error {
	fmt.Println("********Stopping Profile********", debuggerAddress)

	re := regexp.MustCompile(`:(\d+)`)
	match := re.FindStringSubmatch(debuggerAddress)
	if len(match) < 2 {
		return fmt.Errorf("failed to extract port from wsUrl: %s", debuggerAddress)
	}

	port := match[1]
	fmt.Println("Extracted port:", port)

	pid, err := GetProcessPID(port)
	if err != nil {
		return fmt.Errorf("failed to get PID: %v", err)
	}

	// 终止进程
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	fmt.Println("Browser process killed successfully!")
	return nil
}

func GetProcessPID(port string) (int, error) {
	cmd := exec.Command("cmd", "/C", "netstat -ano | findstr LISTENING | findstr :"+port)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get process info: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var foundPID int

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue //  避免解析空行
		}

		//  通常 PID 是最后一个字段
		pid, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			foundPID = pid
			fmt.Println(" Found Browser PID:", pid)
			break //  一旦找到，直接返回
		}
	}

	if foundPID == 0 {
		return 0, fmt.Errorf("no PID found for port %s", port)
	}

	return foundPID, nil
}

func fetchAllProfiles(gl *GoLogin, limit int) ([]map[string]interface{}, error) {
	fmt.Printf("正在进入 fetchAllProfiles()\n")

	all := []map[string]interface{}{}
	page := 1
	seenIDs := make(map[string]bool)

	for {
		url := fmt.Sprintf("https://api.gologin.com/browser/v2?limit=%d&page=%d", limit, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+gl.Token)

		respBody, err := DoRequestWithCheck(req)
		if err != nil {
			return nil, fmt.Errorf("fetch error @ page %d: %v", page, err)
		}

		var raw struct {
			Profiles []map[string]interface{} `json:"profiles"`
		}
		if err := json.Unmarshal(respBody, &raw); err != nil {
			return nil, fmt.Errorf("decode error: %v", err)
		}

		newCount := 0
		for _, p := range raw.Profiles {
			if id, ok := p["id"].(string); ok && !seenIDs[id] {
				seenIDs[id] = true
				all = append(all, p)
				newCount++
			}
		}

		fmt.Printf("page=%d 本页拉取 %d 条，其中新增 %d 条，累计 %d 条\n", page, len(raw.Profiles), newCount, len(all))

		if len(raw.Profiles) < limit || newCount == 0 {
			break
		}

		page++
	}

	return all, nil
}

func CleanupExcessProfiles(gl *GoLogin, maxProfiles int, whitelist map[string]bool, traceID string) (string, error) {
	fmt.Printf(" [%s] 正在进入 CleanupExcessProfiles()\n", traceID)
	profiles, err := fetchAllProfiles(gl, 30)
	if err != nil {
		return "", fmt.Errorf("failed to fetch profiles: %v", err)
	}

	fmt.Printf(" 已获取全部 profiles，共 %d 条\n", len(profiles))

	// 🧹 剔除白名单
	var eligible []map[string]interface{}
	for _, p := range profiles {
		if id, ok := p["id"].(string); ok && !whitelist[id] {
			eligible = append(eligible, p)
		}
	}

	fmt.Printf(" [traceID: %s] 当前可清理 profile 数量: %d（总共 %d）\n", traceID, len(eligible), len(profiles))

	if len(eligible) <= maxProfiles {
		msg := fmt.Sprintf(" 当前 profile 数量 %d 未超过上限 %d，无需清理", len(eligible), maxProfiles)
		logger.LogInfo(traceID, msg)
		return msg, nil
	}

	// 排序：createdAt 越早越靠前
	sort.Slice(eligible, func(i, j int) bool {
		ti, _ := eligible[i]["createdAt"].(string)
		tj, _ := eligible[j]["createdAt"].(string)
		return ti < tj
	})

	toDelete := eligible
	var deleted []string

	for _, p := range toDelete {
		id, _ := p["id"].(string)
		name, _ := p["name"].(string)

		fmt.Println("********Deleting Profile********", id)

		if err := DeleteProfile(gl, id); err != nil {
			msg := fmt.Sprintf(" 删除失败: %s (%s): %v", id, name, err)
			logger.LogError(traceID, msg)
		} else {
			msg := fmt.Sprintf(" 删除成功: %s (%s)", id, name)
			logger.LogInfo(traceID, msg)
			deleted = append(deleted, msg)
		}
	}

	summary := fmt.Sprintf(" 共清理 profile %d 个：\n%s", len(deleted), strings.Join(deleted, "\n"))
	return summary, nil

}

// 删除 Profile
func DeleteProfile(gl *GoLogin, id string) error {
	fmt.Println("********Deleting Profile********", id)
	url := "https://api.gologin.com/browser/" + id
	var isSuccess bool
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Printf("Failed to create DELETE request (Attempt %d): %v\n", i+1, err)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+gl.Token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed to execute DELETE request (Attempt %d): %v\n", i+1, err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			log.Println("Error close response:", err)
		}

		if resp.StatusCode == 204 {
			isSuccess = true
			break //成功后跳出循环
		} else {
			log.Printf("Delete attempt %d failed: status %d\n", i+1, resp.StatusCode)
		}
		time.Sleep(2 * time.Second)

	}
	if !isSuccess {
		return fmt.Errorf("delete failed after 5 attempts")
	}

	fmt.Println("Profile deleted successfully")
	return nil
}

func QuitProfile(gl *GoLogin, profileID string, debuggerAddress string) error {

	var errs []string

	if err := StopProfile(debuggerAddress); err != nil {
		errs = append(errs, fmt.Sprintf("Stop error: %v", err))
	}
	if err := DeleteProfile(gl, profileID); err != nil {
		errs = append(errs, fmt.Sprintf("Delete error: %v", err))
	}
	CleanCache(profileID)

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	fmt.Println("Profile quit successfully")
	return nil
}
