package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	keys       string
	value      string
	initJson   string
	outputPath string
	outputType string
	version    string
	flag       bool
	// staticDir string
)

var RootCmd = &cobra.Command{
	Use:   "gson",
	Short: "go-json-structure[BETA]",
	Long:  `使用go语言，完成一个简单的json解析器命令行工具，支持对json的序列化和反序列化，提供简易的交互UI，同时能对外能提供稳定、安全、统一的api接口`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func init() {
	RootCmd.PersistentFlags().StringVar(&version, "version", "1.0.0", "go-json-structure[BETA]")
	RootCmd.PersistentFlags().BoolVar(&flag, "debug", false, "start with debug mode")

	RootCmd.PersistentFlags().StringVar(&initJson, "json", "", "自定义json数据源(需要满足格式要求)")
	RootCmd.PersistentFlags().StringVar(&keys, "keys", "/", "查询的内容/xxxx/xxx/xxx的形式")
	RootCmd.PersistentFlags().StringVar(&value, "value", "", "增和改设置的json值")
	RootCmd.PersistentFlags().StringVar(&outputType, "output", "json", "输出格式[json|yaml|xml]")
	RootCmd.PersistentFlags().StringVar(&outputPath, "outputPath", "output.txt", "输出地址，默认为output.txt")
	// RootCmd.PersistentFlags().StringVar(&staticDir, "static", "./public/dist", "静态文件夹路径")
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(ServerCmd)

	// =====添加命令行服务
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(InsertCmd)
	RootCmd.AddCommand(UpdateCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(SaveCmd)
}

// OutAlistInit 暴露用于外部启动server的函数
func OutAlistInit() {
	var (
		cmd  *cobra.Command
		args []string
	)
	ServerCmd.Run(cmd, args)
}
