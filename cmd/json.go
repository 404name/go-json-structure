package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/404name/go-json-structure/service"
	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var InsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "添加内容",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := service.UpdateOrInsertValue(getKeys(), value, true)
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf(res)
		}
		os.Exit(0)
	},
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "删除内容",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := service.DeleteValue(getKeys())
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf(res)
		}
		os.Exit(0)
	},
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "修改内容",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := service.UpdateOrInsertValue(getKeys(), value, true)

		if err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf(res)
		}
		os.Exit(0)
	},
}
var SaveCmd = &cobra.Command{
	Use:   "save",
	Short: "保存内容",
	Run: func(cmd *cobra.Command, args []string) {

		if err := service.Save(outputPath, outputType); err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf("已保存到目前目录下：" + outputPath)
		}
		os.Exit(0)
	},
}
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "获取内容",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := service.GetValue(getKeys(), outputType)
		if err != nil {
			fmt.Printf(err.Error())
		} else {
			fmt.Printf(res)
		}
		os.Exit(0)
	},
}

// VersionCmd represents the version command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(service.InitJson(initJson))
		os.Exit(0)
	},
}

func getKeys() []string {
	result := strings.Split(keys, "/")

	var pathArray []string
	for _, part := range result {
		if part != "" {
			pathArray = append(pathArray, part)
		}
	}
	return result
}
