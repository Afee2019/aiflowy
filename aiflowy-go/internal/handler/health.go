package handler

import (
	"runtime"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/config"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

var startTime = time.Now()

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns basic health status
func (h *HealthHandler) Health(c echo.Context) error {
	return response.OK(c, map[string]interface{}{
		"status": "ok",
	})
}

// HealthDetail returns detailed health status
func (h *HealthHandler) HealthDetail(c echo.Context) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Test database connection
	dbStatus := "ok"
	dbInfo, err := repository.TestQuery()
	if err != nil {
		dbStatus = "error: " + err.Error()
		dbInfo = nil
	}

	return response.OK(c, map[string]interface{}{
		"status":      "ok",
		"version":     "1.0.0",
		"environment": config.GetEnv(),
		"uptime":      time.Since(startTime).String(),
		"runtime": map[string]interface{}{
			"go_version":    runtime.Version(),
			"num_goroutine": runtime.NumGoroutine(),
			"num_cpu":       runtime.NumCPU(),
			"os":            runtime.GOOS,
			"arch":          runtime.GOARCH,
		},
		"memory": map[string]interface{}{
			"alloc_mb":       m.Alloc / 1024 / 1024,
			"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
			"sys_mb":         m.Sys / 1024 / 1024,
			"num_gc":         m.NumGC,
		},
		"database": map[string]interface{}{
			"status": dbStatus,
			"info":   dbInfo,
		},
	})
}

// TestHandler handles test endpoints for Stage 1 verification
type TestHandler struct{}

// NewTestHandler creates a new TestHandler
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

// TestError tests error handling
func (h *TestHandler) TestError(c echo.Context) error {
	return apierrors.New(1001, "这是一个测试错误")
}

// TestPanic tests panic recovery
func (h *TestHandler) TestPanic(c echo.Context) error {
	panic("这是一个测试 panic")
}

// TestSnowflake tests snowflake ID generation
func (h *TestHandler) TestSnowflake(c echo.Context) error {
	ids := make([]string, 5)
	for i := 0; i < 5; i++ {
		id, err := snowflake.GenerateIDString()
		if err != nil {
			return apierrors.Wrap(err, "生成 ID 失败")
		}
		ids[i] = id
	}

	// Parse one ID to show components
	id, _ := snowflake.GenerateID()
	timestamp, datacenterID, workerID, sequence := snowflake.ParseID(id)

	return response.OK(c, map[string]interface{}{
		"generated_ids": ids,
		"parsed_example": map[string]interface{}{
			"id":            id,
			"timestamp":     timestamp.Format(time.RFC3339),
			"datacenter_id": datacenterID,
			"worker_id":     workerID,
			"sequence":      sequence,
		},
	})
}

// TestConfig tests configuration loading
func (h *TestHandler) TestConfig(c echo.Context) error {
	cfg := config.Get()
	return response.OK(c, map[string]interface{}{
		"environment":   config.GetEnv(),
		"is_production": config.IsProduction(),
		"server": map[string]interface{}{
			"port": cfg.Server.Port,
			"host": cfg.Server.Host,
			"mode": cfg.Server.Mode,
		},
		"database": map[string]interface{}{
			"host":     cfg.Database.Host,
			"port":     cfg.Database.Port,
			"database": cfg.Database.Database,
		},
	})
}
