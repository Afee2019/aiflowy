package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/config"
	"github.com/aiflowy/aiflowy-go/internal/middleware"
	"github.com/aiflowy/aiflowy-go/internal/repository"
	"github.com/aiflowy/aiflowy-go/internal/router"
	"github.com/aiflowy/aiflowy-go/internal/service/tool/builtin"
	"github.com/aiflowy/aiflowy-go/pkg/logger"
	"github.com/aiflowy/aiflowy-go/pkg/metrics"
	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
	"go.uber.org/zap"
)

var (
	configPath string
	version    = "1.0.0"
	buildTime  = "unknown"
	commitSHA  = "unknown"
)

func init() {
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	err = logger.Init(&logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting AIFlowy Go Backend",
		zap.String("version", version),
		zap.String("build_time", buildTime),
		zap.String("commit_sha", commitSHA),
		zap.String("environment", config.GetEnv()),
	)

	// Initialize Prometheus metrics
	metrics.Init(version, runtime.Version(), buildTime)
	logger.Info("Prometheus metrics initialized")

	// Initialize snowflake ID generator
	err = snowflake.Init(cfg.Snowflake.WorkerID, cfg.Snowflake.DatacenterID)
	if err != nil {
		logger.Fatal("Failed to initialize snowflake", zap.Error(err))
	}
	logger.Info("Snowflake ID generator initialized",
		zap.Int64("worker_id", cfg.Snowflake.WorkerID),
		zap.Int64("datacenter_id", cfg.Snowflake.DatacenterID),
	)

	// Initialize database
	err = repository.InitDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer repository.CloseDB()

	// Register builtin tools
	err = builtin.RegisterAll()
	if err != nil {
		logger.Fatal("Failed to register builtin tools", zap.Error(err))
	}
	logger.Info("Builtin tools registered",
		zap.Int("count", len(builtin.GetBuiltinTools())),
	)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup middleware
	middleware.SetupMiddleware(e)

	// Setup routes
	router.SetupRoutes(e)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	go func() {
		logger.Info("Server starting",
			zap.String("address", addr),
			zap.String("mode", cfg.Server.Mode),
		)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Print startup banner
	printBanner(addr, config.GetEnv())

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func printBanner(addr string, env string) {
	banner := `
    _    ___ _____ _
   / \  |_ _|  ___| | _____      ___   _
  / _ \  | || |_  | |/ _ \ \ /\ / / | | |
 / ___ \ | ||  _| | | (_) \ V  V /| |_| |
/_/   \_\___|_|   |_|\___/ \_/\_/  \__, |
                                   |___/
    AIFlowy Go Backend v%s (%s)

    Server running at: http://%s
    Health check: http://%s/health
    Test endpoints: http://%s/test/*

`
	fmt.Printf(banner, version, env, addr, addr, addr)
}
