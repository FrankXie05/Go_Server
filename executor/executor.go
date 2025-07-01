package executor

import (
	"EMU_server/gologin"
	"EMU_server/scripts"
	"EMU_server/supervisor"
	"EMU_server/types"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RunTask(request types.TaskRequest) (string, int) {
	traceID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	scriptUUID := uuid.New().String()
	fmt.Printf(" [%s] 执行任务 UUID: %s\n", traceID, scriptUUID)
	fmt.Printf("TrackingLink %s:", request.TrackingLink)
	gl := gologin.NewGoLogin()
	config := gologin.NewProfileConfig(scriptUUID, request.Proxy, request.Navigator)

	profileID, err := gologin.CreateProfile(gl, config)
	if err != nil {
		return "create profile failed: " + err.Error(), http.StatusInternalServerError
	}

	debuggerAddress, err := gologin.StartProfile(gl, profileID)
	if err != nil {
		_ = gologin.DeleteProfile(gl, profileID)
		return "start profile failed: " + err.Error(), http.StatusInternalServerError
	}

	// ⏱ 注册任务进 TaskSupervisor
	supervisor.GlobalSupervisor.Register(&supervisor.RunningTask{
		TraceID:      traceID,
		ProfileID:    profileID,
		DebugAddress: debuggerAddress,
		StartTime:    time.Now(),
		Status:       supervisor.StatusRunning,
		Timeout:      90 * time.Second,
		CleanupFunc: func() {
			gologin.QuitProfile(gl, profileID, debuggerAddress)
		},
	})

	// 👇 deferred 版本清理（用于正常流程退出）
	defer func() {
		fmt.Printf(" [%s] 正常退出，清理 profile\n", traceID)
		gologin.QuitProfile(gl, profileID, debuggerAddress)
		supervisor.GlobalSupervisor.MarkDone(traceID)
	}()

	if strings.ToUpper(request.Type) == "CLICK" {
		request.Script = `print("200");sys.exit(0)`
	}

	scriptPath, err := scripts.CreateScript(request.Type, request.TrackingLink, scriptUUID, request.Script)
	if err != nil {
		return "create script failed: " + err.Error(), http.StatusInternalServerError
	}
	defer os.Remove(scriptPath)

	output, err := scripts.RunPythonScript(
		request.Type,
		scriptPath,
		scriptUUID,
		request.TrackingLink,
		debuggerAddress,
		request.SoiMetaData,
	)

	if err != nil {
		return fmt.Sprintf("run script failed: %v\n", err), http.StatusInternalServerError
	}

	return output, http.StatusOK
}
