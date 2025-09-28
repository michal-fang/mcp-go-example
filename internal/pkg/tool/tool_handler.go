// Package tool tool_handler.go
package tool

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// Handler 定义工具处理器的接口
type Handler interface {
	Name() string
	Description() string
	Schema() mcp.Tool
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// BaseTool 基础工具结构体，实现公共功能
type BaseTool struct {
	name        string
	description string
	schema      mcp.Tool
}

func NewBaseTool(name string, description string, schema mcp.Tool) BaseTool {
	return BaseTool{
		name:        name,
		description: description,
		schema:      schema,
	}
}

func (b *BaseTool) Name() string {
	return b.name
}

func (b *BaseTool) Description() string {
	return b.description
}

func (b *BaseTool) Schema() mcp.Tool {
	return b.schema
}
