package scripts

import (
	"EMU_server/types"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 任务脚本存储路径
var scriptDir = "./scripts"
var templatePath = "./scripts/template.py"

// 创建 Python 任务脚本
func CreateScript(ty string, trackingLink string, uuid string, scriptsCode string) (string, error) {
	// 根据类型选择子目录
	subDir := strings.ToUpper(ty)
	scriptSubDir := filepath.Join(scriptDir, subDir)
	// 确保子目录存在
	if err := os.MkdirAll(scriptSubDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create subdir: %v", err)
	}
	filename := fmt.Sprintf("%s.py", uuid)
	scriptPath := filepath.Join(scriptSubDir, filename)

	// 如果脚本已存在，直接返回路径
	if _, err := os.Stat(scriptPath); err == nil {
		return scriptPath, nil
	}

	// 读取 Python 模板
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %v", err)
	}

	// 替换 `{tracking_link}`
	scriptContent := strings.Replace(string(templateData), "{tracking_link}", trackingLink, -1)
	scriptContent = strings.ReplaceAll(scriptContent, "{scripts_code}", scriptsCode)

	// 写入新任务脚本
	err = os.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create script: %v", err)
	}

	return scriptPath, nil
}

func RunPythonScript(ty string, scriptPath string, uuid string, trackingLink string, debuggerAddress string, soiMeta types.SoiMetaData) (string, error) {
	args := []string{scriptPath, uuid, trackingLink, debuggerAddress}

	if !isEmptySoiMeta(soiMeta) {
		metaJSON, err := json.Marshal(soiMeta)
		if err != nil {
			return "", fmt.Errorf("failed to marshal SoiMetaData: %v", err)
		}
		args = append(args, string(metaJSON))
	}

	cmd := exec.Command("python", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("exit code: %v, stderr: %s", err, string(output))
	}
	return string(output), nil
}

// isEmptySoiMeta checks if soiMeta is a zero value.
func isEmptySoiMeta(meta types.SoiMetaData) bool {
	var zero types.SoiMetaData
	return meta == zero
}
