// Package manager
package manager

import (
	"mcp-go-tutorials/tool"

	"github.com/mark3labs/mcp-go/server"
)

// Manager  工具管理器
type Manager struct {
	tools []tool.Handler
}

// NewToolManager 创建工具管理器
func NewToolManager() *Manager {
	return &Manager{
		tools: make([]tool.Handler, 0),
	}
}

// RegisterTool 注册工具
func (tm *Manager) RegisterTool(tool tool.Handler) {
	tm.tools = append(tm.tools, tool)
}

// RegisterAllTools 注册所有工具到MCP服务器
func (tm *Manager) RegisterAllTools(s *server.MCPServer) {
	for _, handler := range tm.tools {
		s.AddTool(handler.Schema(), handler.Handle)
	}
}

// GetTools 获取所有工具
func (tm *Manager) GetTools() []tool.Handler {
	return tm.tools
}
