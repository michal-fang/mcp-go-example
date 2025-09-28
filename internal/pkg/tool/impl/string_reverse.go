package impl

import (
	"context"
	"mcp-go-tutorials/internal/pkg/tool"

	"github.com/mark3labs/mcp-go/mcp"
)

type StringReverseTool struct {
	tool.BaseTool
}

// NewStringReverseTool 创建字符串反转工具
func NewStringReverseTool() tool.Handler {
	reverseTool := mcp.NewTool("reverse_string",
		mcp.WithDescription("Reverse a string"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The text to reverse"),
		),
	)
	return &StringReverseTool{
		BaseTool: tool.NewBaseTool(
			"reverse_string",
			"Reverse a string",
			reverseTool),
	}
}

func (s StringReverseTool) Handle(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return mcp.NewToolResultText(string(runes)), nil
}
