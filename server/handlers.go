package server

import (
	"EMU_server/executor"
	"EMU_server/logger"
	"EMU_server/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("trace-%d", time.Now().UnixNano())
		}
		c.Set("traceId", traceID)
		c.Writer.Header().Set("X-Request-ID", traceID)
		c.Next()
	}
}

func RunTasks(c *gin.Context) {
	traceID := c.GetString("traceId")

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf(" Panic: %v", r)
			logger.LogError(traceID, msg)
			RespondJSON(c, http.StatusInternalServerError, "fail", "Server internal error", nil)
		}
	}()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.LogError(traceID, "读取请求体失败: "+err.Error())
		RespondJSON(c, http.StatusBadRequest, "fail", "Invalid body", nil)
		return
	}

	var requests []types.TaskRequest
	if err := json.Unmarshal(body, &requests); err != nil {
		var single types.TaskRequest
		if err := json.Unmarshal(body, &single); err != nil {
			logger.LogError(traceID, "JSON 解析失败: "+err.Error())
			RespondJSON(c, http.StatusBadRequest, "fail", "Invalid format", nil)
			return
		}
		requests = []types.TaskRequest{single}
	}

	var eg errgroup.Group
	result := make([]string, len(requests))
	codeMap := make([]int, len(requests))

	for i, req := range requests {
		i, req := i, req
		eg.Go(func() error {
			msg, code := executor.RunTask(req)
			result[i] = msg
			codeMap[i] = code
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		logger.LogError(traceID, "任务调度失败: "+err.Error())
		RespondJSON(c, http.StatusInternalServerError, "fail", "Task execution error", nil)
		return
	}

	var failList []string
	worstCode := 0
	for i, code := range codeMap {
		if code != http.StatusOK {
			failList = append(failList, fmt.Sprintf("Task %d error: %s", i, result[i]))
			if code > worstCode {
				worstCode = code
			}
		}
	}

	if len(failList) > 0 {
		msg := strings.Join(failList, "\n")
		logger.LogError(traceID, msg)
		RespondJSON(c, http.StatusMultiStatus, "partial_fail", msg, gin.H{"code": worstCode})
		return
	}

	logger.LogInfo(traceID, "任务全部成功")
	RespondJSON(c, http.StatusOK, "success", "All tasks OK", nil)
}
