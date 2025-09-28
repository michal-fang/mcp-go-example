# mcp-go-example
使用github.com/mark3labs/mcp-go实现mcp server示例
=======
## 运行方式
```shell
# 进入入口
cd cmd

# 使用默认配置
go run main.go

# 使用命令行参数
go run main.go --mode sse --port 8082 --log-level debug

# 使用配置文件
go run main.go --config config.yaml

# 短参数形式
go run main.go -m sse -p 8082
```
