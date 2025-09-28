package app

import (
	"mcp-go-tutorials/pkg/log"
	"mcp-go-tutorials/pkg/version/verflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TransportMode string

const (
	StdioMode TransportMode = "stdio"
	SSEMode   TransportMode = "sse"
	HTTPMode  TransportMode = "streamableHttp"
)

type Config struct {
	Mode     string `mapstructure:"mode"`
	Port     string `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
	GinMode  string `mapstructure:"gin_mode"`
}

var (
	cfgFile string
	cfg     Config
)

func NewAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "mcp-server",
		Short:        "MCP Server with multiple transport modes",
		Long:         `MCP Server supporting stdio, sse, and streamableHttp transport modes`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果 `--version=true`，则打印版本并退出
			verflag.PrintAndExitIfRequested()
			options := log.NewOptions()
			log.Init(options)
			runServer()
			return nil
		},
	}
	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")
	cmd.PersistentFlags().StringP("mode", "m", "streamableHttp", "Transport mode: stdio, sse, http")
	cmd.PersistentFlags().StringP("port", "p", "8081", "Port for HTTP/SSE server")
	cmd.PersistentFlags().String("log-level", "info", "Log level: debug, info, warn, error")
	cmd.PersistentFlags().String("gin-mode", "release", "Gin mode: debug, release, test")

	// 绑定 Viper
	_ = viper.BindPFlag("mode", cmd.PersistentFlags().Lookup("mode"))
	_ = viper.BindPFlag("port", cmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("gin_mode", cmd.PersistentFlags().Lookup("gin-mode"))
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}
