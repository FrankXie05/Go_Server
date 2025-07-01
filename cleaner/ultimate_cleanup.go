package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 清理指定目录下所有内容
func clearDirectoryContents(path string) {
	fmt.Printf("\n[INFO] Cleaning contents of: %s\n", path)

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read: %s → %v\n", path, err)
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		err := os.RemoveAll(fullPath)
		if err != nil {
			fmt.Printf("[FAILED ] %s - %v\n", fullPath, err)
		} else {
			fmt.Printf("[DELETED] %s\n", fullPath)
		}
	}
}

// 无限循环，每 N 分钟清理一次指定路径列表
func StartLocalCleanupForever(interval time.Duration, folders []string) {
	go func() {
		for {
			fmt.Printf(" Local cleaner starting @ %s\n", time.Now().Format(time.RFC3339))
			for _, folder := range folders {
				clearDirectoryContents(folder)
			}
			fmt.Printf(" Cleanup cycle complete. Next in %v...\n", interval)
			time.Sleep(interval)
		}
	}()
}
