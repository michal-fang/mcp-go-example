package app

import (
	"errors"
	"fmt"
	"mcp-go-tutorials/internal/pkg/tool/impl"
	"mcp-go-tutorials/internal/pkg/tool/manager"
	"mcp-go-tutorials/pkg/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func runServer() {
	// 初始化工具管理器
	toolManager := manager.NewToolManager()
	toolManager.RegisterTool(impl.NewCalculatorTool())
	toolManager.RegisterTool(impl.NewStringReverseTool())

	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"MCP Server with multiple transport modes, only support for calculate, string reverse and so on.",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	toolManager.RegisterAllTools(s)
	log.Infof("Starting MCP server in %s mode on port %s", cfg.Mode, cfg.Port)

	switch TransportMode(cfg.Mode) {
	case StdioMode:
		startStdioServer(s)
	case SSEMode:
		startSSEServer(s)
	case HTTPMode:
		startHTTPServer(s)
	default:
		log.Fatalf("Unknown mode: %s", cfg.Mode)
	}
}

func startStdioServer(s *server.MCPServer) {
	fmt.Println("Starting MCP server in stdio mode...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Stdio server error: %v", err)
	}
}

func startSSEServer(s *server.MCPServer) {
	// 使用 Gin 框架
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// 创建 SSE 处理器
	sseHandler := server.NewSSEServer(s)

	// 注册路由
	router.GET("/sse", gin.WrapH(sseHandler))
	router.GET("/health", healthCheckHandler)

	// 优雅关闭支持
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Infof("Starting SSE srv on http://localhost:%s/sse", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("SSE srv error: %v", err)
	}
}

func startHTTPServer(s *server.MCPServer) {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// 创建 HTTP 处理器
	httpHandler := server.NewStreamableHTTPServer(s)

	// 注册路由
	router.POST("/mcp", gin.WrapH(httpHandler))
	router.GET("/mcp", gin.WrapH(httpHandler)) // 支持 GET 请求
	router.GET("/health", healthCheckHandler)
	router.GET("/metrics", metricsHandler)

	// API 文档路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "MCP Server",
			"version": "1.0.0",
			"mode":    cfg.Mode,
			"endpoints": gin.H{
				"mcp":     "/mcp",
				"health":  "/health",
				"metrics": "/metrics",
			},
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Infof("Starting HTTP srv on http://localhost:%s/mcp", cfg.Port)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP srv error: %v", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "MCP Server",
		"mode":    cfg.Mode,
	})
}

func metricsHandler(c *gin.Context) {
	// 这里可以添加应用指标
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"metrics": gin.H{
			"active_connections": 0, // 可以添加实际指标
			"requests_served":    0,
		},
	})
}

func initConfig() {
	setDefaultValue()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("/app")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	// read in environment variables that match
	viper.AutomaticEnv()
	// 读取环境变量的前缀小写,将自动转变为大写。
	viper.SetEnvPrefix("MCP")

	//将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.GinMode)
	log.Infoln("Using config file", "file", viper.ConfigFileUsed())
}

func setDefaultValue() {
	// 设置默认值
	viper.SetDefault("mode", "streamableHttp")
	viper.SetDefault("port", "8081")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("gin_mode", "release")
	//设置日志默认值
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.enableCaller", true)
	viper.SetDefault("log.disableStdout", false)
	viper.SetDefault("log.maxSize", 10)
	viper.SetDefault("log.maxBackups", 3)
	viper.SetDefault("log.maxAge", 7)
	viper.SetDefault("log.timeFormat", "human")
}
