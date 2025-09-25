// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"mcp-go-tutorials/tool/impl"
	"mcp-go-tutorials/tool/manager"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type TransportMode string

const (
	StdioMode TransportMode = "stdio"
	SSEMode   TransportMode = "sse"
	HTTPMode  TransportMode = "streamableHttp"
)

func main() {
	mode := flag.String("mode", "streamableHttp", "Transport mode: stdio, sse, http")
	port := flag.String("port", "8081", "Port for HTTP/SSE server")
	flag.Parse()
	// 创建工具管理器并注册工具
	toolManager := manager.NewToolManager()
	toolManager.RegisterTool(impl.NewCalculatorTool())
	// 可以轻松添加更多工具
	toolManager.RegisterTool(impl.NewStringReverseTool())
	// 创建MCP服务器
	s := server.NewMCPServer(
		"Common Tool",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	// 注册所有工具到服务器
	toolManager.RegisterAllTools(s)
	// 根据模式启动服务器
	switch TransportMode(*mode) {
	case StdioMode:
		startStdioServer(s)
	case SSEMode:
		startSSEServer(s, *port)
	case HTTPMode:
		startHTTPServer(s, *port)
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func startStdioServer(s *server.MCPServer) {
	fmt.Println("Starting MCP server in stdio mode...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Stdio server error: %v", err)
	}
}

func startSSEServer(s *server.MCPServer, port string) {
	handler := server.NewSSEServer(s)

	http.Handle("/sse", handler)
	http.Handle("/health", healthCheckHandler())

	addr := ":" + port
	log.Printf("Starting SSE server on http://localhost%s/sse", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("SSE server error: %v", err)
	}
}

func startHTTPServer(s *server.MCPServer, port string) {
	// 创建HTTP处理器（如果需要普通的HTTP接口）
	handler := server.NewStreamableHTTPServer(s)
	http.Handle("/mcp", handler)
	http.Handle("/health", healthCheckHandler())

	addr := ":" + port
	log.Printf("Starting HTTP server on http://localhost%s/mcp", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func healthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "healthy"}`))
	})
}
