package supervisor

import (
	"fmt"
	"sync"
	"time"
)

var GlobalSupervisor = NewSupervisor()

type TaskStatus string

const (
	StatusRunning TaskStatus = "running"
	StatusDone    TaskStatus = "done"
	StatusTimeout TaskStatus = "timeout"
	StatusFailed  TaskStatus = "failed"
)

// 跟踪正在运行的任务
type RunningTask struct {
	TraceID      string
	ProfileID    string
	DebugAddress string
	StartTime    time.Time
	Status       TaskStatus
	Timeout      time.Duration
	CleanupFunc  func()
}

// 管理所有任务
type TaskSupervisor struct {
	tasks map[string]*RunningTask
	mu    sync.RWMutex
}

// 创建一个新的 Supervisor 实例
func NewSupervisor() *TaskSupervisor {
	return &TaskSupervisor{
		tasks: make(map[string]*RunningTask),
	}
}

// 注册任务
func (ts *TaskSupervisor) Register(task *RunningTask) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tasks[task.TraceID] = task
	fmt.Printf(" 注册任务 [%s], Profile: %s\n", task.TraceID, task.ProfileID)

	// 启动 watchdog 监控超时
	go ts.watchTask(task)
}

// 标记任务完成
func (ts *TaskSupervisor) MarkDone(traceID string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if task, ok := ts.tasks[traceID]; ok {
		task.Status = StatusDone
		fmt.Printf("任务完成 [%s]\n", traceID)
	}
}

// 内部超时监控
func (ts *TaskSupervisor) watchTask(task *RunningTask) {
	<-time.After(task.Timeout)

	ts.mu.Lock()
	defer ts.mu.Unlock()
	if task.Status == StatusRunning {
		fmt.Printf("任务超时 [%s], 开始清理 Profile: %s\n", task.TraceID, task.ProfileID)
		task.Status = StatusTimeout
		if task.CleanupFunc != nil {
			task.CleanupFunc()
		}
	}
}

// 获取当前任务快照（调试用）
func (ts *TaskSupervisor) ListActive() []RunningTask {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	list := []RunningTask{}
	for _, task := range ts.tasks {
		if task.Status == StatusRunning {
			list = append(list, *task)
		}
	}
	return list
}
