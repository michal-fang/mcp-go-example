// Package impl calculator.go
package impl

import (
	"context"
	"fmt"
	"mcp-go-tutorials/tool"

	"github.com/mark3labs/mcp-go/mcp"
)

// CalculatorTool 计算器工具实现
type CalculatorTool struct {
	tool.BaseTool
}

// NewCalculatorTool 创建新的计算器工具实例
func NewCalculatorTool() tool.Handler {
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	return &CalculatorTool{
		BaseTool: tool.NewBaseTool(
			"calculate",
			"Perform basic arithmetic operations",
			calculatorTool),
	}
}

// Handle 处理计算器工具请求
func (c *CalculatorTool) Handle(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	x, err := request.RequireFloat("x")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	y, err := request.RequireFloat("y")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var result float64
	switch op {
	case "add":
		result = x + y
	case "subtract":
		result = x - y
	case "multiply":
		result = x * y
	case "divide":
		if y == 0 {
			return mcp.NewToolResultError("cannot divide by zero"), nil
		}
		result = x / y
	default:
		return mcp.NewToolResultError("unknown operation: " + op), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
}
