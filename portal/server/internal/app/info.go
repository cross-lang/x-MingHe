package app

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// InfoResponse 服务信息响应
type InfoResponse struct {
	Version         string `json:"version"`          // 版本号
	BuildTime       string `json:"build_time"`       // 构建时间
	GitCommit       string `json:"git_commit"`       // Git 提交哈希
	GoVersion       string `json:"go_version"`       // Go 版本
	OS              string `json:"os"`               // 操作系统
	Arch            string `json:"arch"`             // 系统架构
	StartTime       string `json:"start_time"`       // 服务启动时间
	Uptime          string `json:"uptime"`           // 运行时长
	Goroutines      int    `json:"goroutines"`       // 协程数量
	MemoryStats     string `json:"memory_stats"`     // 内存使用情况
	DependencyCount int    `json:"dependency_count"` // 依赖包数量
}

var (
	startTime      = time.Now()
	// 构建时通过 ldflags 注入的变量
	Version   = "v1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// InfoHandler 服务信息处理器
// @Description 获取服务当前状态和系统信息
// @Produce json
// @Success 200 {object} ginx.Response
// @Router /info [get]
func InfoHandler(ctx *gin.Context) {
	uptime := time.Since(startTime)

	// 获取内存统计
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryInfo := fmt.Sprintf("Alloc=%vMB, TotalAlloc=%vMB, Sys=%vMB, NumGC=%v",
		bToMb(memStats.Alloc),
		bToMb(memStats.TotalAlloc),
		bToMb(memStats.Sys),
		memStats.NumGC,
	)

	response := InfoResponse{
		Version:     Version,
		BuildTime:   BuildTime,
		GitCommit:   GitCommit,
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		StartTime:   startTime.Format("2006-01-02 15:04:05"),
		Uptime:      uptime.String(),
		Goroutines:  runtime.NumGoroutine(),
		MemoryStats: memoryInfo,
	}

	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
}

// bToMb 将字节转换为 MB
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
